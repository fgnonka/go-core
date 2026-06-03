package main

import (
	"go-core/internal/data"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	// Create a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize our in-memory store
	store := data.NewStore()

	// Create an instance of our application with the store dependency
	app := &application{store: store}

	mux := http.NewServeMux()

	// Register explicit CRUD API endpoints using pattern matching placeholders
	mux.HandleFunc("GET /api/v1/entities", app.handleListTeams)
	mux.HandleFunc("POST /api/v1/entities/{id}/tip", app.handleProcessTip)
	mux.HandleFunc("POST /api/v1/entities/{id}/withdraw", app.handleProcessWithdrawal)

	// Create a server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("starting wallet platform server", "addr", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
