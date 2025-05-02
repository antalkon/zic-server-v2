package ws

import (
	"context"
	"time"
)

func (t *Tunnel) startPing(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	timeout := 60 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			if now.Sub(t.LastPong) > timeout {
				t.Conn.Close()
				return
			}

			t.Conn.WriteJSON(Message{
				Version:   "1.0",
				Type:      "ping",
				ID:        "ping-" + t.ID,
				From:      "server",
				Timestamp: now,
				Payload: map[string]interface{}{
					"timestamp": now.Format(time.RFC3339),
					"sequence":  time.Now().Unix(),
				},
			})
		}
	}
}
