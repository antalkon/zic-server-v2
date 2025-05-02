package ws

import "time"

type Message struct {
	Version   string      `json:"version"`
	Type      string      `json:"type"`
	ID        string      `json:"id"`
	From      string      `json:"from"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}
