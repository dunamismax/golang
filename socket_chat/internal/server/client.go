// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: client.go
// author: dunamismax
// version: 1.0.0
// date: 17-06-2025
// github: <https://github.com/dunamismax>
// description: Defines the Client type, which represents a single user's connection and handles concurrent I/O.

package server

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"time"
)

const (
	// writeWait is the time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// pongWait is the time allowed to read the next pong message from the peer.
	// For this TCP implementation, any read from the client will reset the deadline.
	pongWait = 60 * time.Second
	// pingPeriod is the interval for sending pings to the peer. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// maxMessageSize is the maximum message size allowed from a peer.
	maxMessageSize = 1024
)

// Client is a middleman between the TCP connection and the hub.
// It encapsulates all connection-specific state and logic.
type Client struct {
	// hub is a reference to the central hub that manages message broadcasting.
	hub *Hub

	// conn is the underlying network connection for the client.
	conn net.Conn

	// send is a buffered channel for outbound messages. Messages placed in this
	// channel are sent to the client's connection by the writePump.
	send chan []byte

	// nickname is the identifier for the client in the chat.
	nickname string
}

// newClient creates a new Client instance.
func newClient(hub *Hub, conn net.Conn, nickname string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		nickname: nickname,
	}
}

// readPump pumps messages from the TCP connection to the hub.
//
// This method runs in its own goroutine for each client, allowing for concurrent
// reads without blocking the rest of the server. It reads messages from the
// client's connection, formats them with the client's nickname, and forwards
// them to the hub's broadcast channel.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	reader := bufio.NewReader(c.conn)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				slog.Warn("Client connection timed out", "nickname", c.nickname, "remote_addr", c.conn.RemoteAddr())
			} else {
				slog.Info("Client disconnected", "nickname", c.nickname, "remote_addr", c.conn.RemoteAddr(), "error", err.Error())
			}
			break
		}

		// Reset the read deadline each time we successfully read data.
		c.conn.SetReadDeadline(time.Now().Add(pongWait))

		// Format the message and send it to the hub to be broadcast.
		// The newline character is trimmed from the end.
		message := []byte(fmt.Sprintf("[%s]: %s", c.nickname, line))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the TCP connection.
//
// This method also runs in its own goroutine for each client. It waits for
// messages on the client's `send` channel and writes them to the connection.
// It also sends periodic pings to the client to ensure the connection is alive.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.Write([]byte("Server is closing the connection.\n"))
				return
			}

			if _, err := c.conn.Write(message); err != nil {
				slog.Error("Failed to write message to client", "nickname", c.nickname, "error", err)
				return
			}

		case <-ticker.C:
			// Send a ping message (a simple newline) to keep the connection alive.
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if _, err := c.conn.Write([]byte("\n")); err != nil {
				slog.Error("Failed to send ping to client", "nickname", c.nickname, "error", err)
				return
			}
		}
	}
}
