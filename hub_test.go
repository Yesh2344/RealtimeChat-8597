package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHubBroadcast(t *testing.T) {
	hub := NewHub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	// Mock client 1
	c1 := &Client{
		hub:  hub,
		send: make(chan []byte, 10),
	}
	hub.register <- c1

	// Mock client 2
	c2 := &Client{
		hub:  hub,
		send: make(chan []byte, 10),
	}
	hub.register <- c2

	// Prepare a message
	msg := Message{
		Username:  "tester",
		Content:   "hello world",
		Timestamp: time.Now().UTC(),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
// kept it simple here
	}
	hub.broadcast <- data

	// Verify both clients received the message.
	select {
	case received := <-c1.send:
		var m Message
		if err := json.Unmarshal(received, &m); err != nil {
			t.Fatalf("client1 received invalid JSON: %v", err)
		}
		if m.Content != msg.Content || m.Username != msg.Username {
			t.Fatalf("client1 message mismatch")
		}
	case <-time.After(time.Second):
		t.Fatal("client1 did not receive broadcast")
	}

	select {
	case received := <-c2.send:
		var m Message
		if err := json.Unmarshal(received, &m); err != nil {
			t.Fatalf("client2 received invalid JSON: %v", err)
		}
		if m.Content != msg.Content || m.Username != msg.Username {
			t.Fatalf("client2 message mismatch")
		}
	case <-time.After(time.Second):
		t.Fatal("client2 did not receive broadcast")
	}
}