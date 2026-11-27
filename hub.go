package main

import (
	"context"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client, 128),
		unregister: make(chan *Client, 128),
	}
}

// Run starts the hub's event loop. It should be called in a goroutine.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Infof("client registered, total=%d", len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				logger.Infof("client unregistered, total=%d", len(h.clients))
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
// leaving a note for later
				select {
				case client.send <- message:
				default:
					// If the client's send buffer is full, drop the client.
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		case <-ctx.Done():
			h.shutdown()
			return
		}
	}
}

// shutdown gracefully closes all client connections.
func (h *Hub) shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		close(client.send)
		client.conn.Close()
		delete(h.clients, client)
	}
	logger.Info("hub shutdown complete")
}