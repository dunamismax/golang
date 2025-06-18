# Gopher v2.0 Demo: Go, Alpine.js, & Tailwind CSS

This project is a minimal, production-ready web application that demonstrates the Gopher v2.0 architecture. It features a Go backend serving a simple HTML page styled with Tailwind CSS and made interactive with Alpine.js.

All frontend assets are embedded into the final, single Go binary for easy deployment.

## Core Architecture

- **Backend**: Go (`net/http`, `chi` router)
- **Frontend**: Alpine.js for interactivity, Tailwind CSS for styling.
- **Templating**: Go's `html/template` package.
- **Assets**: Go's `embed` package for bundling.
- **Logging**: Go's `slog` for structured JSON logging.

## Prerequisites

You must have the following tools installed on your system. This project is optimized for **macOS (ARM64)**.

1. **Go**: Version 1.21 or newer.
2. **Node.js & npm**: Required for the Tailwind CSS build step.

## Project Structure

```sh
/gopher-v2-demo
├── go.mod
├── go.sum
├── main.go               # Main application: server, routing, handlers.
├── package.json          # Frontend build dependencies (Tailwind).
├── tailwind.config.js    # Tailwind CSS configuration.
├── README.md             # This file.
│
├── static/               # Directory for assets to be embedded.
│   ├── alpine.v3.min.js  # Alpine.js library.
│   └── output.css        # Compiled Tailwind CSS.
│
└── web/
    ├── assets/
    │   └── input.css     # Source Tailwind CSS file.
    └── templates/
        └── index.gohtml  # Go HTML template.
```

## Build and Run Instructions

Follow these steps exactly to build and run the application.

### Step 1: Set Up Project Directory

First, create the project directory and the necessary subdirectories.

```sh
mkdir -p gopher-v2-demo/web/assets gopher-v2-demo/web/templates gopher-v2-demo/static
cd gopher-v2-demo
```

### Step 2: Create Project Files

Create the following files inside the `gopher-v2-demo` directory. Copy and paste the content provided for each file.

- `main.go`
- `web/templates/index.gohtml`
- `web/assets/input.css`
- `tailwind.config.js`
- `package.json`

### Step 3: Initialize Go Module

Initialize the Go module. This will create the `go.mod` file.

```sh
go mod init gopher-v2-demo
```

### Step 4: Install Go Dependencies

Get the `chi` router dependency. This will also create the `go.sum` file.

```sh
go get github.com/go-chi/chi/v5
```

### Step 5: Download Alpine.js

Download the Alpine.js library directly into the `static` directory.

```sh
curl -o static/alpine.v3.min.js https://cdn.jsdelivr.net/npm/alpinejs@3.14.1/dist/cdn.min.js
```

### Step 6: Build Frontend Assets

Install the frontend dependency (Tailwind CSS) and compile `input.css` into `static/output.css`.

```sh
# Install npm packages defined in package.json
npm install

# Run the build script defined in package.json to compile the CSS
npm run build
```

You should now see a generated `static/output.css` file.

### Step 7: Build and Run the Go Application

Build the final, single binary. The `go build` command will compile the Go code and embed the `static` directory's contents into the executable.

```sh
# Build the binary
go build -o gopher-demo .

# Run the application
./gopher-demo
```

You will see a JSON log message in your terminal:

```json
{
  "time": "2025-06-18T15:15:00.000Z",
  "level": "INFO",
  "msg": "Starting server",
  "addr": ":3000"
}
```

Open your web browser and navigate to **[http://localhost:3000](http://localhost:3000)**. You will see the running application.

---
