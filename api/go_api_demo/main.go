// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: main.go
// author:   dunamismax
// version:  2.0.0
// date:     06-19-2025
// github:   <https://github.com/dunamismax>
// description: The ultimate single-file, production-ready Go HTTP API server
// demonstration. This version includes CRUD operations, repository pattern,
// middleware, struct validation, and concurrent in-memory storage.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
)

// =============================================================================
// 1. DOMAIN MODELS & VALIDATION
// =============================================================================

// User represents the data model for a user in our system.
// Field tags for `json` and `validate` are used for serialization and validation.
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Email     string    `json:"email" validate:"required,email"`
}

// = a new validator instance.
var validate = validator.New()

// =============================================================================
// 2. REPOSITORY PATTERN (DATA LAYER)
// =============================================================================

// UserRepository defines the interface for user data storage.
// This allows us to decouple the application from the specific database implementation.
// All methods accept a context for cancellation and timeout propagation.
type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, id string, user User) (User, error)
	Delete(ctx context.Context, id string) error
}

// InMemoryUserRepository is a thread-safe, in-memory implementation of UserRepository.
// It uses a sync.RWMutex to handle concurrent read/write operations safely.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]User
}

// NewInMemoryUserRepository creates and returns a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]User),
	}
}

// Create adds a new user to the in-memory store.
func (r *InMemoryUserRepository) Create(ctx context.Context, user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Simulate context cancellation check
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	if _, exists := r.users[user.ID]; exists {
		return User{}, fmt.Errorf("user with ID %s already exists", user.ID)
	}

	r.users[user.ID] = user
	return user, nil
}

// GetByID retrieves a user by their ID from the in-memory store.
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	user, exists := r.users[id]
	if !exists {
		return User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetAll retrieves all users from the store.
func (r *InMemoryUserRepository) GetAll(ctx context.Context) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	allUsers := make([]User, 0, len(r.users))
	for _, user := range r.users {
		allUsers = append(allUsers, user)
	}
	return allUsers, nil
}

// Update modifies an existing user in the store.
func (r *InMemoryUserRepository) Update(ctx context.Context, id string, user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	if _, exists := r.users[id]; !exists {
		return User{}, fmt.Errorf("user not found")
	}

	user.ID = id // Ensure the ID remains the same
	r.users[id] = user
	return user, nil
}

// Delete removes a user from the store.
func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(r.users, id)
	return nil
}

// =============================================================================
// 3. APPLICATION & DEPENDENCY INJECTION
// =============================================================================

// Config holds all configuration for the application.
// Values are read from environment variables.
type Config struct {
	Port string
}

// application is the central struct holding all application-wide dependencies,
// such as the logger and data models (repositories).
type application struct {
	config Config
	logger *slog.Logger
	users  UserRepository
}

// =============================================================================
// 4. HTTP HANDLERS, HELPERS, & ROUTING
// =============================================================================

// jsonResponse is a generic structure for sending JSON responses.
type jsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// errorResponse is a generic structure for sending JSON error messages.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON is a helper for sending JSON-formatted responses.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.logger.Error("failed to write JSON response", "error", err)
	}
}

// writeError is a helper for sending JSON-formatted error responses.
func (app *application) writeError(w http.ResponseWriter, status int, message string) {
	app.writeJSON(w, status, errorResponse{Error: message})
}

// readJSON is a helper that decodes JSON from the request body and validates it.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Limit request body size to 1MB.
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Prevent unknown fields in the request body.

	if err := dec.Decode(dst); err != nil {
		// Handle various decoding errors with user-friendly messages.
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at character %d", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return errors.New("body must not be larger than 1MB")
		default:
			return err
		}
	}

	// Check for a second JSON value in the body.
	if dec.More() {
		return errors.New("body must only contain a single JSON value")
	}

	// Perform validation on the decoded struct.
	if err := validate.Struct(dst); err != nil {
		// This can be further enhanced to return more detailed validation errors.
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// loggingMiddleware logs details of each incoming HTTP request.
func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		app.logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
		)
	})
}

// routes sets up all the application routes and applies middleware.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Registering user-related handlers under a versioned API path.
	mux.HandleFunc("POST /api/v1/users", app.createUserHandler)
	mux.HandleFunc("GET /api/v1/users", app.getAllUsersHandler)
	mux.HandleFunc("GET /api/v1/users/{id}", app.getUserHandler)
	mux.HandleFunc("PUT /api/v1/users/{id}", app.updateUserHandler)
	mux.HandleFunc("DELETE /api/v1/users/{id}", app.deleteUserHandler)

	return app.loggingMiddleware(mux)
}

// --- CRUD Handlers ---

// createUserHandler handles the creation of a new user.
// POST /api/v1/users
//
//	curl -X POST -H "Content-Type: application/json" \
//	 -d '{"name": "dunamismax", "email": "dev@example.com"}' \
//	 http://localhost:8080/api/v1/users
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name" validate:"required,min=2,max=100"`
		Email string `json:"email" validate:"required,email"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user := User{
		ID:        fmt.Sprintf("user_%d", time.Now().UnixNano()), // Simple unique ID
		CreatedAt: time.Now(),
		Name:      input.Name,
		Email:     input.Email,
	}

	createdUser, err := app.users.Create(r.Context(), user)
	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	app.writeJSON(w, http.StatusCreated, jsonResponse{
		Status:  "success",
		Message: "User created successfully",
		Data:    createdUser,
	})
}

// getUserHandler retrieves a user by their ID from the path parameter.
// GET /api/v1/users/{id}
// curl http://localhost:8080/api/v1/users/{id}
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user, err := app.users.GetByID(r.Context(), id)
	if err != nil {
		app.writeError(w, http.StatusNotFound, "user not found")
		return
	}

	app.writeJSON(w, http.StatusOK, jsonResponse{
		Status: "success",
		Data:   user,
	})
}

// getAllUsersHandler retrieves all users.
// GET /api/v1/users
// curl http://localhost:8080/api/v1/users
func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.users.GetAll(r.Context())
	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "could not retrieve users")
		return
	}

	app.writeJSON(w, http.StatusOK, jsonResponse{
		Status: "success",
		Data:   users,
	})
}

// updateUserHandler handles updating an existing user.
// PUT /api/v1/users/{id}
//
//	curl -X PUT -H "Content-Type: application/json" \
//	 -d '{"name": "dunamismax_v2", "email": "dev_v2@example.com"}' \
//	 http://localhost:8080/api/v1/users/{id}
func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input struct {
		Name  string `json:"name" validate:"required,min=2,max=100"`
		Email string `json:"email" validate:"required,email"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch existing user to update. In a real app, you might only update certain fields.
	existingUser, err := app.users.GetByID(r.Context(), id)
	if err != nil {
		app.writeError(w, http.StatusNotFound, "user not found")
		return
	}

	existingUser.Name = input.Name
	existingUser.Email = input.Email

	updatedUser, err := app.users.Update(r.Context(), id, existingUser)
	if err != nil {
		app.writeError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	app.writeJSON(w, http.StatusOK, jsonResponse{
		Status:  "success",
		Message: "User updated successfully",
		Data:    updatedUser,
	})
}

// deleteUserHandler handles deleting a user.
// DELETE /api/v1/users/{id}
// curl -X DELETE http://localhost:8080/api/v1/users/{id}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := app.users.Delete(r.Context(), id); err != nil {
		app.writeError(w, http.StatusNotFound, "user not found")
		return
	}

	app.writeJSON(w, http.StatusOK, jsonResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

// =============================================================================
// 5. MAIN APPLICATION ENTRYPOINT
// =============================================================================

func main() {
	// 1. Initialize logger.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 2. Load configuration.
	cfg := Config{
		Port: os.Getenv("API_PORT"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// 3. Initialize dependencies (database repository).
	userRepo := NewInMemoryUserRepository()

	// 4. Create the main application struct with all dependencies.
	app := &application{
		config: cfg,
		logger: logger,
		users:  userRepo,
	}

	// 5. Configure the HTTP server.
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      app.routes(), // Use the new router
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 6. Run the server in a goroutine for graceful shutdown.
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		logger.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			shutdownError <- err
		}

		logger.Info("server shutdown complete")
		shutdownError <- nil
	}()

	logger.Info("server starting", "address", srv.Addr)

	// Start the server. If it fails for reasons other than a clean shutdown,
	// log the error.
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed to start", "error", err)
		os.Exit(1)
	}

	// Block until the shutdown process is complete.
	if err := <-shutdownError; err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
}
