package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type WeeklyGoalRepository struct {
	db *sql.DB
}

func NewWeeklyGoalRepository(db *sql.DB) *WeeklyGoalRepository {
	return &WeeklyGoalRepository{db: db}
}

func (r *WeeklyGoalRepository) Upsert(ctx context.Context, g *domain.WeeklyGoal) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO weekly_goals (user_id, week_start, target_min)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, week_start)
		DO UPDATE SET target_min = EXCLUDED.target_min
		`, g.UserID, g.WeekStart, g.TargetMin)
	return err
}

func (r *WeeklyGoalRepository) FindByUserAndWeek(ctx context.Context, userID string, weekStart time.Time) (*domain.WeeklyGoal, error) {
	var g domain.WeeklyGoal
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, week_start, target_min
		FROM weekly_goals
		WHERE user_id = $1 AND week_start = $2
		`, userID, weekStart).Scan(&g.UserID, &g.WeekStart, &g.TargetMin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrWeeklyGoalNotFound
		}
		return nil, err
	}
	return &g, nil
}
