// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: game.go
// author:   dunamismax
// version:  1.0.1
// date:     06-18-2025
// github:   <https://github.com/dunamismax>
// description: Core game logic and state management for the Snake game.
package game

import (
	"fmt"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont" // Correct font import
)

const (
	// GridSize defines the size of each cell in the game grid.
	GridSize = 10
	// TicksPerMove defines how many game ticks pass before the snake moves one step.
	// This controls the game speed.
	TicksPerMove = 5
)

var (
	backgroundColor = color.RGBA{50, 50, 50, 255}
	gameOverColor   = color.RGBA{255, 0, 0, 255}
	scoreColor      = color.RGBA{255, 255, 255, 255}
)

// Game holds the entire state of the game.
// It implements the ebiten.Game interface.
type Game struct {
	screenWidth  int
	screenHeight int
	gridWidth    int
	gridHeight   int
	snake        *Snake
	food         *Food
	score        int
	tickCounter  int
	gameOver     bool
	needsReset   bool
}

// NewGame is the constructor for a new Game object.
func NewGame(width, height int) (*Game, error) {
	g := &Game{
		screenWidth:  width,
		screenHeight: height,
		gridWidth:    width / GridSize,
		gridHeight:   height / GridSize,
	}
	g.reset()
	return g, nil
}

// reset initializes or resets the game to its starting state.
func (g *Game) reset() {
	slog.Info("Resetting game state")
	g.snake = NewSnake(g.gridWidth/2, g.gridHeight/2)
	g.food = NewFood(g.gridWidth, g.gridHeight, g.snake.body)
	g.score = 0
	g.tickCounter = 0
	g.gameOver = false
	g.needsReset = false
}

// Update progresses the game state by one tick.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.needsReset = true
		}
		// Must reset on the next tick to avoid race conditions with input handling.
		if g.needsReset {
			g.reset()
		}
		return nil
	}

	// Handle player input
	g.snake.HandleInput()

	g.tickCounter++
	if g.tickCounter < TicksPerMove {
		return nil
	}
	g.tickCounter = 0

	// Move snake and check for game events
	g.snake.Move()

	// Check for wall collision
	if g.snake.CollidesWithWall(g.gridWidth, g.gridHeight) {
		g.endGame()
		return nil
	}

	// Check for self collision
	if g.snake.CollidesWithSelf() {
		g.endGame()
		return nil
	}

	// Check for eating food
	if g.snake.Head().X == g.food.pos.X && g.snake.Head().Y == g.food.pos.Y {
		g.score++
		g.snake.Grow()
		g.food = NewFood(g.gridWidth, g.gridHeight, g.snake.body)
		slog.Info("Food eaten", "new_score", g.score)
	}

	return nil
}

// endGame sets the game over state.
func (g *Game) endGame() {
	slog.Warn("Game Over", "final_score", g.score)
	g.gameOver = true
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	g.snake.Draw(screen)
	g.food.Draw(screen)

	// Use the imported basicfont.Face7x13 for rendering text. [6]
	fontFace := basicfont.Face7x13

	// Draw Score
	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, fontFace, 10, 20, scoreColor)

	if g.gameOver {
		msg := "GAME OVER\n(Press R to Restart)"
		bounds := text.BoundString(fontFace, msg)
		x := (g.screenWidth - bounds.Dx()) / 2
		y := (g.screenHeight - bounds.Dy()) / 2
		text.Draw(screen, msg, fontFace, x, y, gameOverColor)
	}
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}
