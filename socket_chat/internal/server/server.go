// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: server.go
// author: dunamismax
// version: 1.0.0
// date: 17-06-2025
// github: <https://github.com/dunamismax>
// description: Defines the main TCP server for handling client connections.

package server

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
)

// Server handles incoming TCP connections and manages the chat hub.
type Server struct {
	address string
	hub     *Hub
}

// NewServer creates a new Server instance.
func NewServer(address string) *Server {
	return &Server{
		address: address,
		hub:     NewHub(),
	}
}

// Start initializes the server, starts the hub, and listens for connections.
// This is the primary entry point for the chat server's lifecycle.
func (s *Server) Start() error {
	// The hub is the core of our concurrency model. It must be running in the
	// background to process registrations, unregistrations, and broadcasts.
	go s.hub.Run()

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start listener on %s: %w", s.address, err)
	}
	defer listener.Close()

	slog.Info("Starting chat server", "address", s.address)

	// The main server loop. It blocks here, waiting for new connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			// In a real-world scenario, you might want to introduce a backoff
			// delay here instead of continuing immediately.
			continue
		}

		// We handle each new connection in its own goroutine to allow the
		// server to immediately go back to accepting more connections.
		go s.handleConnection(conn)
	}
}

// handleConnection manages a single client connection from acceptance until termination.
func (s *Server) handleConnection(conn net.Conn) {
	slog.Info("New client connected", "remote_addr", conn.RemoteAddr())

	// Prompt for nickname.
	_, err := conn.Write([]byte("Enter your nickname: "))
	if err != nil {
		slog.Error("Failed to write nickname prompt", "remote_addr", conn.RemoteAddr(), "error", err)
		conn.Close()
		return
	}

	// Read nickname from the client.
	reader := bufio.NewReader(conn)
	nickname, err := reader.ReadString('\n')
	if err != nil {
		slog.Error("Failed to read nickname", "remote_addr", conn.RemoteAddr(), "error", err)
		conn.Close()
		return
	}
	// Trim whitespace and newline characters from the nickname.
	nickname = nickname[:len(nickname)-1]

	// Create a new client instance for this connection.
	client := newClient(s.hub, conn, nickname)

	// Register the new client with the hub. This is a channel send, which
	// will be processed by the hub's single goroutine.
	s.hub.register <- client

	// Announce the new user to the chat.
	joinMsg := fmt.Sprintf("%s has joined the chat.\n", nickname)
	s.hub.broadcast <- []byte(joinMsg)

	// Start the I/O pumps. These run in their own goroutines, allowing for
	// concurrent reading and writing for this client. The readPump will handle
	// unregistration and connection closing when it exits.
	go client.writePump()
	go client.readPump()
}
