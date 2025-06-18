// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: hub.go
// author: dunamismax
// version: 1.0.0
// date: 17-06-2025
// github: <https://github.com/dunamismax>
// description: Defines the concurrent hub that manages chat clients and message broadcasting.

package server

import "log/slog"

// Hub maintains the set of active clients and broadcasts messages to the
// clients. It is the central component of the concurrent chat server.
type Hub struct {
	// clients is the set of all currently connected clients. The map keys are
	// pointers to Client objects, and the values are boolean `true`.
	clients map[*Client]bool

	// broadcast is the channel for incoming messages from the clients. Messages
	// sent to this channel will be forwarded to all connected clients.
	broadcast chan []byte

	// register is the channel for new clients wishing to register with the hub.
	register chan *Client

	// unregister is the channel for clients that wish to unregister from the hub.
	unregister chan *Client
}

// NewHub creates and returns a new Hub instance, ready to be run.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop. This method should be run as a goroutine.
// The hub's single goroutine design ensures that access to the `clients` map
// is serialized, preventing race conditions without the need for mutexes.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// A new client has connected. Add it to the clients map.
			h.clients[client] = true
			slog.Info("Client registered", "nickname", client.nickname, "remote_addr", client.conn.RemoteAddr())

		case client := <-h.unregister:
			// A client has disconnected. Check if it exists, then remove it and
			// close its send channel to signal the writePump to exit.
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				slog.Info("Client unregistered", "nickname", client.nickname, "remote_addr", client.conn.RemoteAddr())
			}

		case message := <-h.broadcast:
			// A message has been received from a client. Iterate through all
			// connected clients and send the message to their respective send channels.
			for client := range h.clients {
				select {
				case client.send <- message:
					// Message was successfully queued for the client.
				default:
					// The client's send channel is full. This indicates a slow
					// consumer. Close the channel and unregister the client.
					slog.Warn("Client send buffer full. Disconnecting.", "nickname", client.nickname)
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
