// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: main.go
// author: dunamismax
// version: 1.0.0
// date: 17-06-2025
// github: <https://github.com/dunamismax>
// description: The entrypoint for the socket_chat command-line client.

package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	// serverAddress is the address of the chat server to connect to.
	serverAddress = "localhost:8080"
)

func main() {
	// Pillar IV: Structured Logging is mandatory from inception.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Pillar IV: Design for graceful shutdown.
	// Set up a context that is canceled when a SIGINT or SIGTERM is received.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	slog.Info("Connecting to chat server", "address", serverAddress)
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		slog.Error("Failed to connect to server", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	slog.Info("Connection successful. Use Ctrl+C to exit or Ctrl+D to stop sending.")
	fmt.Println("----------------------------------------------------------------")

	// This channel and sync.Once ensure that the main goroutine is notified
	// as soon as either the reading or writing goroutine terminates,
	// preventing a panic from closing the channel twice.
	done := make(chan struct{})
	var once sync.Once
	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	// Pillar I: The Concurrency Mandate.
	// I/O operations are handled in dedicated goroutines to prevent blocking.
	go readFromServer(conn, closeDone)
	go writeToServer(conn, closeDone)

	// Block until either the OS sends a shutdown signal or one of the I/O
	// goroutines finishes (e.g., the server disconnects or stdin is closed).
	select {
	case <-ctx.Done():
		slog.Info("Shutdown signal received.")
	case <-done:
		slog.Info("Connection closed or input stream ended.")
		// Give a moment for final messages to be printed.
		time.Sleep(100 * time.Millisecond)
	}

	slog.Info("Client shutting down.")
	fmt.Println("----------------------------------------------------------------")
}

// readFromServer handles reading all incoming data from the server and printing
// it to standard output. It runs in its own goroutine.
func readFromServer(conn net.Conn, done func()) {
	defer done()
	if _, err := io.Copy(os.Stdout, conn); err != nil {
		// This error is expected when the connection is closed either by the
		// client or server.
		if err != io.EOF {
			slog.Debug("Error reading from server", "error", err)
		}
	}
	slog.Info("Server connection has been closed.")
}

// writeToServer handles reading all user input from standard input and sending
// it to the server. It runs in its own goroutine.
func writeToServer(conn net.Conn, done func()) {
	defer done()
	if _, err := io.Copy(conn, os.Stdin); err != nil {
		// This can happen if the connection is closed while writing.
		slog.Error("Failed to send data to server", "error", err)
	}
	slog.Info("Stopped sending messages. You can still receive messages.")
}
