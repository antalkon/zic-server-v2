package ws

import (
	"backend/internal/transport/service"
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
	payload, ok := m.Payload.(map[string]interface{})
	if !ok {
		log.Println("❌ Неверный формат payload")
		return
	}

	t.ComputerID = payload["computer_id"].(string)

	info := map[string]interface{}{
		"tunnel_id": t.ID,
		"status":    "online",
		"last_seen": time.Now().Format(time.RFC3339),
		"os":        payload["os"],
		"ip":        payload["local_ip"],
	}

	jsonValue, err := json.Marshal(info)
	if err != nil {
		log.Println("❌ Redis marshal error:", err)
		return
	}

	if err := t.Redis.Set("pc:"+t.ComputerID, jsonValue, 2*time.Minute); err != nil {
		log.Println("❌ Redis set error:", err)
	}

	// Ответ клиенту
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
