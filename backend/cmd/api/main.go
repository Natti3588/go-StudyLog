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

	// --- repository層 ---
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	studyLogRepo := repository.NewStudyLogRepository(db)
	userStatsRepo := repository.NewUserStatsRepository(db)
	weeklyGoalRepo := repository.NewWeeklyGoalRepository(db)

	// --- service層 ---
	authService := service.NewAuthService(userRepo, jwtSecret)
	categoryService := service.NewCategoryService(categoryRepo)
	studyLogService := service.NewStudyLogService(studyLogRepo, userStatsRepo, db)
	weeklyGoalService := service.NewWeeklyGoalService(weeklyGoalRepo)
	statsService := service.NewStatsService(studyLogRepo, userStatsRepo, weeklyGoalRepo)

	// --- handler層 ---
	authHandler := handler.NewAuthHandler(authService, secureCookie)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	studyLogHandler := handler.NewStudyLogHandler(studyLogService)
	weeklyGoalHandler := handler.NewWeeklyGoalHandler(weeklyGoalService)
	statsHandler := handler.NewStatsHandler(statsService)

	requireAuth := handler.RequireAuth(jwtSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /signup", authHandler.Signup)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", requireAuth(authHandler.Logout))
	mux.HandleFunc("GET /me", requireAuth(authHandler.Me))

	mux.HandleFunc("GET /categories", requireAuth(categoryHandler.List))
	mux.HandleFunc("POST /categories", requireAuth(categoryHandler.Create))

	mux.HandleFunc("GET /logs", requireAuth(studyLogHandler.List))
	mux.HandleFunc("POST /logs", requireAuth(studyLogHandler.Create))
	mux.HandleFunc("PUT /logs/{id}", requireAuth(studyLogHandler.Update))
	mux.HandleFunc("DELETE /logs/{id}", requireAuth(studyLogHandler.Delete))

	mux.HandleFunc("PUT /goals/weekly", requireAuth(weeklyGoalHandler.SetWeekly))
	mux.HandleFunc("GET /stats/summary", requireAuth(statsHandler.Summary))
	mux.HandleFunc("GET /stats/heatmap", requireAuth(statsHandler.Heatmap))

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:5173"
	}
	corsMiddleware := handler.CORS(corsOrigin)

	slog.Info("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", corsMiddleware(mux)); err != nil {
		slog.Error("server failed to start", "error", err)
	}
}
