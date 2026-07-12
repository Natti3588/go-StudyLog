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

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	secureCookie := os.Getenv("COOKIE_SECURE") == "true"

	authHandler := handler.NewAuthHandler(
		service.NewAuthService(repository.NewUserRepository(db), jwtSecret),
		secureCookie,
	)
	categoryHandler := handler.NewCategoryHandler(
		service.NewCategoryService(repository.NewCategoryRepository(db)),
	)
	studyLogHandler := handler.NewStudyLogHandler(
		service.NewStudyLogService(repository.NewStudyLogRepository(db)),
	)

	requireAuth := handler.RequireAuth(jwtSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", authHandler.Signup)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", requireAuth(authHandler.Logout))

	mux.HandleFunc("GET /categories", requireAuth(categoryHandler.List))
	mux.HandleFunc("POST /categories", requireAuth(categoryHandler.Create))

	mux.HandleFunc("GET /logs", requireAuth(studyLogHandler.List))
	mux.HandleFunc("POST /logs", requireAuth(studyLogHandler.Create))
	mux.HandleFunc("PUT /logs/{id}", requireAuth(studyLogHandler.Update))
	mux.HandleFunc("DELETE /logs/{id}", requireAuth(studyLogHandler.Delete))

	slog.Info("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("server failed to start", "error", err)
	}
}
