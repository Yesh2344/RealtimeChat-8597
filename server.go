package main

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader upgrades HTTP connections to WebSocket.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// In production you should verify the origin.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// serveWs handles WebSocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("WebSocket upgrade error: %v", err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Start read/write pumps.
	go client.writePump()
	go client.readPump()
}

// NewServer creates an http.Server with routes configured.
func NewServer(hub *Hub, addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	// Simple health endpoint.
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	handler := corsMiddleware(mux)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return srv
}