package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // потом можно ограничить
	},
}

// Основной хендлер WebSocket
func HandleTunnel(w http.ResponseWriter, r *http.Request, tunnelID string) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade failed:", err)
		return err
	}

	client := &Connection{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	tunnel := &Tunnel{
		ID:         tunnelID,
		Connection: client,
	}

	go tunnel.HandleMessages()

	log.Printf("✅ Установлено WS соединение для туннеля: %s\n", tunnelID)
	return nil
}
