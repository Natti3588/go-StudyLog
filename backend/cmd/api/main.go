package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		slog.Warn("no .env file loaded", "error", err)
	}

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to open db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]any{})
	})

	slog.Info("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server failed to start", "error", err)
	}
}
