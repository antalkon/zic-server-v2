package zicprotws

import (
	"backend/internal/tunel_service"
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
	Service    *tunel_service.TunelService // 👈 добавляем
}

var ActiveTunnels = make(map[string]*Tunnel)
