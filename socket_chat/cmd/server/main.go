// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: main.go
// author: dunamismax
// version: 1.0.0
// date: 17-06-2025
// github: <https://github.com/dunamismax>
// description: The entrypoint for the socket_chat server application.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dunamismax/golang/socket_chat/internal/server"
)

const (
	// defaultAddress is the default host and port the server will listen on.
	defaultAddress = "localhost:8080"
)

func main() {
	// Pillar IV: Structured Logging is mandatory from inception.
	// Initialize a structured JSON logger and set it as the default.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Pillar I & IV: Design for graceful shutdown.
	// Set up a context that listens for interruption signals (SIGINT, SIGTERM).
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Initialize the core application server.
	chatServer := server.NewServer(defaultAddress)

	// Use a channel to listen for the server's startup error. The server's
	// Start method is blocking, so we run it in a goroutine.
	errChan := make(chan error, 1)
	go func() {
		// The Start method will only return an error if it fails to initialize,
		// for instance, if the port is already in use.
		errChan <- chatServer.Start()
	}()

	slog.Info("Server process initiated. Waiting for signal or fatal error.")

	// Block until a signal is received or the server's Start method returns a fatal error.
	select {
	case err := <-errChan:
		// A fatal error occurred during server startup or operation.
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	case <-ctx.Done():
		// A shutdown signal (SIGINT or SIGTERM) was received.
		// The context being "Done" is our trigger to begin a clean exit.
		slog.Info("Shutdown signal received. Terminating server.")
	}

	// At this point, the program is about to exit. The defer statements in
	// server.Start() will handle closing the listener. Active connections
	// will be terminated as their goroutines exit.
	fmt.Println("\nServer has been shut down.")
	os.Exit(0)
}
