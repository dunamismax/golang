// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: snake.go
// author:   dunamismax
// version:  1.0.0
// date:     06-18-2025
// github:   <https://github.com/dunamismax>
// description: Defines the snake object, its behavior, and rendering.
package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	snakeHeadColor = color.RGBA{0, 255, 127, 255} // Spring Green
	snakeBodyColor = color.RGBA{0, 200, 100, 255}
)

// Direction represents the snake's direction of movement.
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

// Position represents a coordinate on the game grid.
type Position struct {
	X, Y int
}

// Snake represents the player-controlled snake.
type Snake struct {
	body      []Position
	direction Direction
	nextDir   Direction
	growth    int
}

// NewSnake creates a new snake at a starting position.
func NewSnake(startX, startY int) *Snake {
	s := &Snake{
		body:      make([]Position, 3), // Start with 3 segments
		direction: DirRight,
		nextDir:   DirRight,
	}
	s.body[0] = Position{startX, startY}
	s.body[1] = Position{startX - 1, startY}
	s.body[2] = Position{startX - 2, startY}
	return s
}

// HandleInput checks for keyboard input to change the snake's direction.
func (s *Snake) HandleInput() {
	if (inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)) && s.direction != DirDown {
		s.nextDir = DirUp
	} else if (inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)) && s.direction != DirUp {
		s.nextDir = DirDown
	} else if (inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)) && s.direction != DirRight {
		s.nextDir = DirLeft
	} else if (inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)) && s.direction != DirLeft {
		s.nextDir = DirRight
	}
}

// Move updates the snake's position based on its current direction.
func (s *Snake) Move() {
	s.direction = s.nextDir
	head := s.body[0]
	newHead := head

	switch s.direction {
	case DirUp:
		newHead.Y--
	case DirDown:
		newHead.Y++
	case DirLeft:
		newHead.X--
	case DirRight:
		newHead.X++
	}

	// Prepend the new head and manage growth
	s.body = append([]Position{newHead}, s.body...)

	if s.growth > 0 {
		s.growth--
	} else {
		// Remove the tail if not growing
		s.body = s.body[:len(s.body)-1]
	}
}

// Grow marks the snake to grow on its next move.
func (s *Snake) Grow() {
	s.growth++
}

// Head returns the position of the snake's head.
func (s *Snake) Head() Position {
	return s.body[0]
}

// CollidesWithWall checks if the snake's head has hit the game boundaries.
func (s *Snake) CollidesWithWall(gridWidth, gridHeight int) bool {
	head := s.Head()
	return head.X < 0 || head.Y < 0 || head.X >= gridWidth || head.Y >= gridHeight
}

// CollidesWithSelf checks if the snake's head has hit its own body.
func (s *Snake) CollidesWithSelf() bool {
	head := s.Head()
	for i, segment := range s.body {
		if i == 0 { // Skip the head itself
			continue
		}
		if head.X == segment.X && head.Y == segment.Y {
			return true
		}
	}
	return false
}

// Draw renders the snake on the screen.
func (s *Snake) Draw(screen *ebiten.Image) {
	for i, segment := range s.body {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(segment.X*GridSize), float64(segment.Y*GridSize))

		// Create a single pixel image to be scaled, for efficiency
		img := ebiten.NewImage(1, 1)
		if i == 0 {
			img.Fill(snakeHeadColor)
		} else {
			img.Fill(snakeBodyColor)
		}

		// Create a temporary image to draw the rectangle
		rect := ebiten.NewImage(GridSize-1, GridSize-1) // -1 for a grid effect
		if i == 0 {
			rect.Fill(snakeHeadColor)
		} else {
			rect.Fill(snakeBodyColor)
		}
		screen.DrawImage(rect, opts)
	}
}
