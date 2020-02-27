package chat

import (
	"dyc/internal/logger"
)

type Responser interface {
	Bytes() []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[_clientId]*Client

	// Inbound messages from the clients.
	broadcast chan Responser

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Responser),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[_clientId]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			logger.Debugf("广播一条新消息: %s", message.Bytes())
			for id, client := range h.clients {
				select {
				case client.send <- message.Bytes():
				default:
					close(client.send)
					delete(h.clients, id)
				}
			}
		}
	}
}
