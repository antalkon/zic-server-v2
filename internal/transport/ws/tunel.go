package ws

import (
	"backend/internal/transport/service"
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
	Service    *service.TunelService // 👈 добавляем
}

var ActiveTunnels = make(map[string]*Tunnel)
