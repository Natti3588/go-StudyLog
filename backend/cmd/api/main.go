package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/Natti3588/go-StudyLog/backend/internal/handler"
	"github.com/Natti3588/go-StudyLog/backend/internal/repository"
	"github.com/Natti3588/go-StudyLog/backend/internal/service"
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

	categoryHandler := handler.NewCategoryHandler(
		service.NewCategoryService(repository.NewCategoryRepository(db)),
	)
	studyLogHandler := handler.NewStudyLogHandler(
		service.NewStudyLogService(repository.NewStudyLogRepository(db)),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /categories", categoryHandler.List)
	mux.HandleFunc("POST /categories", categoryHandler.Create)

	mux.HandleFunc("GET /logs", studyLogHandler.List)
	mux.HandleFunc("POST /logs", studyLogHandler.Create)
	mux.HandleFunc("PUT /logs/{id}", studyLogHandler.Update)
	mux.HandleFunc("DELETE /logs/{id}", studyLogHandler.Delete)

	slog.Info("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server failed to start", "error", err)
	}
}
