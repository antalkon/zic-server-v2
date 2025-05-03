package zicprotws

import (
	wsmodels "backend/internal/ws_models"
	"context"
	"encoding/json"
	"log"
	"time"
)

func (t *Tunnel) listen(ctx context.Context) {
	defer func() {
		t.Conn.Close()
		delete(ActiveTunnels, t.ID)
		t.Cancel()
		t.Redis.Del("pc:" + t.ComputerID)

		// ⚠️ Вызываем Disconnect, если ComputerID указан
		if t.ComputerID != "" && t.Service != nil {
			err := t.Service.Disconnect(&wsmodels.InitPayload{
				ComputerID: t.ComputerID,
			})
			if err != nil {
				log.Printf("❌ Ошибка при отключении компьютера (%s): %v", t.ComputerID, err)
			} else {
				log.Printf("📴 Компьютер %s успешно отключён через Disconnect", t.ComputerID)
			}
		}

		log.Println("🛑 Соединение закрыто:", t.ID)
	}()

	for {
		_, msg, err := t.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		switch m.Type {
		case "init":
			t.handleInit(m)
		case "pong":
			t.LastPong = time.Now()

			if t.ComputerID != "" {
				info := map[string]interface{}{
					"tunnel_id": t.ID,
					"status":    "online",
					"last_seen": time.Now().Format(time.RFC3339),
				}
				jsonValue, _ := json.Marshal(info)
				t.Redis.Set("pc:"+t.ComputerID, jsonValue, 2*time.Minute)
			}
		default:
			log.Println("Unsupported type:", m.Type)
		}
	}
}
