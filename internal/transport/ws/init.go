package ws

import (
	"backend/internal/transport/service"
	wsmodels "backend/internal/ws_models"
	"backend/pkg/cache"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func InitTunnel(conn *websocket.Conn, tunnelID string, redis *cache.RedisClient, service *service.TunelService) {
	ctx, cancel := context.WithCancel(context.Background())

	t := &Tunnel{
		ID:       tunnelID,
		Conn:     conn,
		Cancel:   cancel,
		Redis:    redis,
		Service:  service,
		LastPong: time.Now(),
	}
	ActiveTunnels[tunnelID] = t

	go t.listen(ctx)
	go t.startPing(ctx)
}

func (t *Tunnel) handleInit(m Message) {
	var payload wsmodels.InitPayload
	raw, err := json.Marshal(m.Payload)
	if err != nil {
		log.Println("❌ Failed to marshal payload:", err)
		t.sendErrorAndClose(m.ID, "Ошибка сериализации payload")
		return
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Println("❌ Invalid payload format:", err)
		t.sendErrorAndClose(m.ID, "Неверный формат payload")
		return
	}

	t.ComputerID = payload.ComputerID
	if err := t.Service.GetTunnelByID(&payload); err != nil {
		log.Printf("❌ Failed to get tunnel by ID: %v\n", err)
		t.sendErrorAndClose(m.ID, "Ошибка инициализации туннеля: "+err.Error())
		return
	}

	info := map[string]interface{}{
		"tunnel_id": t.ID,
		"status":    "online",
		"last_seen": time.Now().Format(time.RFC3339),
		"os":        payload.OS,
		"ip":        payload.LocalIP,
	}

	jsonValue, err := json.Marshal(info)
	if err != nil {
		log.Println("❌ Redis marshal error:", err)
		t.sendErrorAndClose(m.ID, "Ошибка сохранения данных в Redis")
		return
	}

	if err := t.Redis.Set("pc:"+t.ComputerID, jsonValue, 2*time.Minute); err != nil {
		log.Println("❌ Redis set error:", err)
		t.sendErrorAndClose(m.ID, "Ошибка записи в Redis")
		return
	}

	resp := Message{
		Version:   "1.0",
		Type:      "response",
		ID:        m.ID,
		From:      "server",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"status":             "ok",
			"tunnel_id":          t.ID,
			"message":            "Успешно инициализирован",
			"server_time":        time.Now().Format(time.RFC3339),
			"ping_interval":      30000,
			"pong_timeout":       60000,
			"reconnect_attempts": 3,
			"reconnect_delay":    5000,
		},
	}
	t.Conn.WriteJSON(resp)
}

func (t *Tunnel) sendErrorAndClose(id, errMsg string) {
	errorMsg := Message{
		Version:   "1.0",
		Type:      "error",
		ID:        id,
		From:      "server",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"status":  "error",
			"message": errMsg,
		},
	}
	_ = t.Conn.WriteJSON(errorMsg)
	t.Conn.Close()
}
