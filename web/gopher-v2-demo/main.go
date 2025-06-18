// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: main.go
// author:   dunamismax
// version:  1.0.0
// date:     18-06-2025
// github:   <https://github.com/dunamismax>
// description: A lightweight Go web server demonstrating the Gopher v2.0 stack.
// It serves a simple, interactive frontend using chi, Tailwind CSS, and Alpine.js.
// All frontend assets are embedded into the binary.

package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed all:static
var staticFS embed.FS

//go:embed web/templates/index.gohtml
var templateFile string

func main() {
	// 1. Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 2. Setup router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger) // Chi's default logger is fine for this demo
	r.Use(middleware.Recoverer)

	// 3. Define routes and handlers
	// Create a file system for static assets
	staticRoot, err := fs.Sub(staticFS, "static")
	if err != nil {
		logger.Error("failed to create static file system", "error", err)
		os.Exit(1)
	}

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticRoot))))

	// Serve the main page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Define data to pass to the template
		pageData := struct {
			Title   string
			Message string
		}{
			Title:   "Gopher v2.0 Demo",
			Message: "Go + Tailwind + Alpine.js",
		}

		// Parse the embedded template file
		tmpl, err := template.New("index").Parse(templateFile)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logger.Error("could not parse template", "error", err)
			return
		}

		// Execute the template
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, pageData)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logger.Error("could not execute template", "error", err)
		}
	})

	// 4. Start the server
	port := ":3000"
	logger.Info("Starting server", "addr", port)
	if err := http.ListenAndServe(port, r); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
