package ws

import (
	"log"
)

type Tunnel struct {
	ID         string
	Connection *Connection
}

// Обработка сообщений в туннеле
func (t *Tunnel) HandleMessages() {
	for {
		_, msg, err := t.Connection.Conn.ReadMessage()
		if err != nil {
			log.Printf("❌ Ошибка чтения из туннеля [%s]: %v", t.ID, err)
			break
		}
		log.Printf("📨 [%s] Получено: %s", t.ID, string(msg))
		// TODO: обрабатывать сообщение по протоколу
	}
}
