package zicprotws

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func SendCommandDirect(tunnelID string, command string, args map[string]interface{}) error {
	tunnel, ok := ActiveTunnels[tunnelID]
	if !ok {
		return fmt.Errorf("активный туннель с ID %s не найден", tunnelID)
	}

	msg := Message{
		Version:   "1.0",
		Type:      "command",
		ID:        "msg-" + uuid.NewString(),
		From:      "server",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"command": command,
			"args":    args,
		},
	}

	return tunnel.Conn.WriteJSON(msg)
}
