// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: food.go
// author:   dunamismax
// version:  1.0.0
// date:     06-18-2025
// github:   <https://github.com/dunamismax>
// description: Defines the food object, its placement, and rendering.
package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	foodColor = color.RGBA{227, 23, 13, 255} // A shade of red
)

// Food represents a food item for the snake to eat.
type Food struct {
	pos   Position
	image *ebiten.Image
}

// NewFood creates a food item at a random position, avoiding the snake's body.
func NewFood(gridWidth, gridHeight int, snakeBody []Position) *Food {
	// Create a map of snake body positions for quick lookup
	snakePosMap := make(map[Position]bool)
	for _, pos := range snakeBody {
		snakePosMap[pos] = true
	}

	var newPos Position
	for {
		newPos = Position{
			X: rand.Intn(gridWidth),
			Y: rand.Intn(gridHeight),
		}
		// Ensure the new food position is not on the snake
		if !snakePosMap[newPos] {
			break
		}
	}

	// Create a reusable image for the food
	img := ebiten.NewImage(GridSize-1, GridSize-1)
	img.Fill(foodColor)

	return &Food{
		pos:   newPos,
		image: img,
	}
}

// Draw renders the food on the screen.
func (f *Food) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(f.pos.X*GridSize), float64(f.pos.Y*GridSize))
	screen.DrawImage(f.image, opts)
}
