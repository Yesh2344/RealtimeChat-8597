package main

import "time"

// Message represents a chat message exchanged over WebSocket.
type Message struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}