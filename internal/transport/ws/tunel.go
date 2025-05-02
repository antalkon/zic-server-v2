package ws

import (
	"backend/pkg/cache"
	"context"
	"time"

	"github.com/gorilla/websocket"
)

type Tunnel struct {
	ID         string
	ComputerID string
	Conn       *websocket.Conn
	LastPong   time.Time
	Cancel     context.CancelFunc
	Redis      *cache.RedisClient
}

var ActiveTunnels = make(map[string]*Tunnel)
