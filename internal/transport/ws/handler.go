package ws

import (
	"backend/pkg/cache"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleTunnel(redis *cache.RedisClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		tunnelID := c.Param("uuid")

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return err
		}

		InitTunnel(conn, tunnelID, redis)
		return nil
	}
}
