package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Connection wraps websocket.Conn
type Connection struct {
	Conn *websocket.Conn
	Send chan []byte // канал для отправки сообщений
	mu   sync.Mutex  // защита от конкурентной записи
}

// Write безопасная запись
func (c *Connection) Write(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteMessage(websocket.TextMessage, msg)
}
