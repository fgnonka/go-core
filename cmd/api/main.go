package main

import (
	"database/sql"
	"go-core/internal/data"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
)

func main() {
	// Create a logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Open a local file-based SQLite database
	db, err := sql.Open("sqlite3", "wallet.db")
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Set connection pool limits appropriate for local SQLite concurrency
	db.SetMaxOpenConns(1)

	schema := `
	CREATE TABLE IF NOT EXISTS entities (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		wallet_balance REAL NOT NULL
	);
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		entity_id TEXT NOT NULL,
		amount REAL NOT NULL,
		type TEXT NOT NULL,
		timestamp TEXT NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		slog.Error("failed to initialize schema", "error", err)
		os.Exit(1)
	}

	// Seed data if database is totally empty
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM entities").Scan(&count)
	if count == 0 {
		_, _ = db.Exec("INSERT INTO entities (id, name, wallet_balance) VALUES ('asec', 'ASEC Mimosas FC', 800.0)")
		_, _ = db.Exec("INSERT INTO entities (id, name, wallet_balance) VALUES ('afad', 'Academie Fad', 500.0)")
		_, _ = db.Exec("INSERT INTO entities (id, name, wallet_balance) VALUES ('afc', 'Africa Sports FC', 300.0)")
		_, _ = db.Exec("INSERT INTO entities (id, name, wallet_balance) VALUES ('soa', 'SOA FC', 200.0)")
	}

	store := data.NewStore(db)
	// Create an instance of our application with the store dependency
	app := &application{store: store}

	mux := http.NewServeMux()

	// Register explicit CRUD API endpoints using pattern matching placeholders
	mux.HandleFunc("GET /api/v1/entities", app.handleListEntities)
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
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
