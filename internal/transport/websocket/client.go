package websocket

import (
	"encoding/json"
	"log"
	"time"

	"backend/internal/tunnel"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn       *websocket.Conn
	hub        *Hub
	tunnelID   string
	computerID string
	send       chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		if c.tunnelID != "" {
			// Удаляем туннель из Redis при закрытии соединения
			if err := c.hub.redisClient.DeleteTunnel(c.tunnelID); err != nil {
				log.Printf("error deleting tunnel: %v", err)
			}
		}
		c.conn.Close()
	}()

	for {
		// Устанавливаем таймаут чтения в 1 минуту
		c.conn.SetReadDeadline(time.Now().Add(time.Minute))
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Обновляем TTL туннеля при получении любого сообщения
		if c.tunnelID != "" {
			if err := c.hub.redisClient.UpdateTunnelTTL(c.tunnelID); err != nil {
				log.Printf("error updating tunnel TTL: %v", err)
			}
		}

		log.Printf("received message: %s", message)

		// Parse and validate init message
		var initMessage InitMessage
		if err := json.Unmarshal(message, &initMessage); err != nil {
			log.Printf("error parsing init message: %v", err)
			continue
		}

		log.Printf("parsed init message: %+v", initMessage)

		if initMessage.Type != "init" {
			log.Printf("unexpected message type: %s", initMessage.Type)
			continue
		}

		log.Println("init message validated")

		// TODO: validate JWT token

		// TODO: check if computer exists in DB

		// Create new tunnel
		tunnelID := "generated-tunnel-id"
		tunnel := &tunnel.Tunnel{
			ID:         tunnelID,
			ComputerID: initMessage.Payload.ComputerID,
			CreatedAt:  time.Now().Unix(),
		}

		err = c.hub.redisClient.SaveTunnel(tunnel)
		if err != nil {
			log.Printf("error saving tunnel: %v", err)
			return
		}

		log.Printf("tunnel created: %s", tunnelID)

		c.tunnelID = tunnelID
		c.computerID = initMessage.Payload.ComputerID

		// Send response
		response := ResponseMessage{
			Version:   "1.0",
			Type:      "response",
			ID:        initMessage.ID,
			From:      "server",
			Timestamp: time.Now().Format(time.RFC3339),
			Payload: ResponsePayload{
				Status:   "ok",
				TunnelID: tunnelID,
				Message:  "Успешно инициализирован",
			},
		}

		responseBytes, _ := json.Marshal(response)
		c.send <- responseBytes

		log.Printf("response sent: %s", responseBytes)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Send message to tunnel
			err = c.hub.redisClient.PublishToTunnel(c.tunnelID, message)
			if err != nil {
				log.Printf("error publishing to tunnel: %v", err)
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

const writeWait = 10 * time.Second

type InitMessage struct {
	Version   string `json:"version"`
	Type      string `json:"type"`
	ID        string `json:"id"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp"`
	Payload   struct {
		ComputerID string `json:"computer_id"`
		JWT        string `json:"jwt"`
	} `json:"payload"`
}

type ResponseMessage struct {
	Version   string          `json:"version"`
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	From      string          `json:"from"`
	Timestamp string          `json:"timestamp"`
	Payload   ResponsePayload `json:"payload"`
}

type ResponsePayload struct {
	Status   string `json:"status"`
	TunnelID string `json:"tunnel_id"`
	Message  string `json:"message"`
}

func (c *Client) handleInitMessage(message []byte) {
	var initMessage InitMessage
	if err := json.Unmarshal(message, &initMessage); err != nil {
		log.Printf("invalid init message: %v", err)
		c.conn.Close()
		return
	}

	if initMessage.Type != "init" {
		log.Printf("unexpected message type: %s", initMessage.Type)
		c.conn.Close()
		return
	}

	// Генерация ID туннеля
	tunnelID := "tunnel-" + initMessage.Payload.ComputerID

	tun := &tunnel.Tunnel{
		ID:         tunnelID,
		ComputerID: initMessage.Payload.ComputerID,
		CreatedAt:  time.Now().Unix(),
	}

	if err := c.hub.redisClient.SaveTunnel(tun); err != nil {
		log.Printf("failed to save tunnel: %v", err)
		c.conn.Close()
		return
	}

	c.tunnelID = tunnelID
	c.computerID = initMessage.Payload.ComputerID

	// Отправляем ответ
	resp := ResponseMessage{
		Version:   "1.0",
		Type:      "response",
		ID:        initMessage.ID,
		From:      "server",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: ResponsePayload{
			Status:   "ok",
			TunnelID: tunnelID,
			Message:  "Туннель успешно создан",
		},
	}

	respBytes, _ := json.Marshal(resp)
	c.send <- respBytes
}
