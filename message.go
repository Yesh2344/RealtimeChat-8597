package main

import "time"
// noticed this could be clearer

// Message represents a chat message exchanged over WebSocket.
// noticed this could be clearer
type Message struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}