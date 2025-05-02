package handlers

import (
	"backend/internal/transport/ws"
	"backend/internal/tunel_service"
	"backend/pkg/cache"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type TunelHandler struct {
	tunel    *tunel_service.TunelService
	validate *validator.Validate
}

func NewTunelHandler(tunel *tunel_service.TunelService) *TunelHandler {
	return &TunelHandler{
		tunel:    tunel,
		validate: validator.New(),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *TunelHandler) HandleTunnel(redis *cache.RedisClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		tunnelID := c.Param("uuid")

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Println("❌ WebSocket Upgrade error:", err)
			return err
		}

		ws.InitTunnel(conn, tunnelID, redis, h.tunel) // 👈 передаём сервис
		return nil
	}
}
