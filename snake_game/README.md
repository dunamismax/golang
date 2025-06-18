# Gopher Demo: Go Snake Game with Ebitengine

This project is a minimal, production-ready game that demonstrates the Gopher development philosophy. It features a complete Snake game built entirely in Go using the high-performance 2D game library, Ebitengine.

The application is self-contained and compiles to a single, statically-linked binary with no external assets required for deployment. This project resides within the `snake_game/` directory of the parent `golang` repository.

## Core Architecture

- **Engine**: Go with Ebitengine (`github.com/hajimehoshi/ebiten/v2`) for rendering, input, and the main game loop.
- **Dependencies**: The Go standard library is used for core functionality.
- **Font Rendering**: `golang.org/x/image/font/basicfont` is used for displaying the score and game-over text.
- **Logging**: Go's `slog` for structured JSON logging to standard output.
- **Build System**: The standard Go toolchain (`go build`, `go run`, `go mod`).

## Prerequisites

You must have the following tools installed on your system. This project is optimized for **macOS (ARM64)**.

1. **Go**: Version 1.21 or newer.

## Local Development Workflow

Follow these steps exactly to build and run the application after cloning the parent repository.

### Step 1: Clone the Repository

First, clone the entire `golang` repository to your local machine.

```sh
git clone https://github.com/dunamismax/golang.git
cd golang
```

### Step 2: Navigate to Project Directory

All subsequent commands must be run from within this specific project's directory.

```sh
cd snake_game
```

### Step 3: Install Go Dependencies

Ensure all Go module dependencies are downloaded and consistent with the `go.mod` file.

```sh
go mod tidy
```

### Step 4: Run the Application

Execute the following command to compile and run the game.

```sh
go run .
```

Alternatively, you can build a distributable binary:

```sh
# Build the binary named "snake"
go build -o snake .

# Run the compiled application
./snake
```

### Step 5: Play the Game

A game window will appear on your screen. You will see JSON log messages in your terminal confirming the game has started.

**Controls:**

- **Arrow Keys / W, A, S, D**: Change the snake's direction.
- **Q**: Quit the game at any time.
- **R**: Restart the game after a "Game Over".

---

## Project Structure Reference

For reference, here is the file structure of this specific game within the monorepo:

```sh
/snake_game
├── go.mod                # Go module definition.
├── go.sum                # Go module checksums.
├── main.go               # Application entry point, initializes Ebitengine.
├── README.md             # This file.
│
└── game/                 # Package for all core game logic.
    ├── game.go           # Main game state, update loop, and drawing logic.
    ├── snake.go          # Snake object, movement, and collision logic.
    └── food.go           # Food object and placement logic.
```
