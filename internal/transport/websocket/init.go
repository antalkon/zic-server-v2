package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func InitWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		computerID: "",
		tunnelID:   "",
	}

	client.hub.register <- client

	go client.writePump()

	// Чтение init-сообщения напрямую из WebSocket
	conn.SetReadDeadline(time.Now().Add(25 * time.Second))
	_, initMsgBytes, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error reading init message: %v", err)
		} else if err == websocket.ErrReadLimit {
			log.Println("init message read timeout")
		}
		conn.Close()
		return
	}

	go client.readPump() // можно запускать после init

	// Обрабатываем init сообщение
	client.handleInitMessage(initMsgBytes)
}
