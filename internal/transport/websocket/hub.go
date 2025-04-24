package websocket

import (
	"backend/internal/tunnel"
	"context"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	redisClient *tunnel.RedisClient
}

func NewHub(redisClient *tunnel.RedisClient) *Hub {
	return &Hub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		redisClient: redisClient,
	}
}

func (h *Hub) listenForMessages() {
	pubsub := h.redisClient.Client.Subscribe(context.Background())
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Printf("error receiving message from Redis: %v", err)
			continue
		}

		tunnelID := msg.Channel

		for client := range h.clients {
			if client.tunnelID == tunnelID {
				select {
				case client.send <- []byte(msg.Payload):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}

	go h.listenForMessages()
}

// ... existing code ...
