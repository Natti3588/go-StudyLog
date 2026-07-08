package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func main() {
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
