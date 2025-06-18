// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: main.go
// author:   dunamismax
// version:  1.0.0
// date:     06-18-2025
// github:   <https://github.com/dunamismax>
// description: The main entry point for the Snake game application.
package main

import (
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"snake.game/m/game"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	// 1. Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Create a new game instance
	g, err := game.NewGame(screenWidth, screenHeight)
	if err != nil {
		slog.Error("failed to initialize game", "error", err)
		os.Exit(1)
	}

	// 3. Configure and run the game window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake by dunamismax")

	slog.Info("Starting game loop")
	if err := ebiten.RunGame(g); err != nil {
		slog.Error("game loop exited with error", "error", err)
		os.Exit(1)
	}
	slog.Info("Game exited cleanly")
}
