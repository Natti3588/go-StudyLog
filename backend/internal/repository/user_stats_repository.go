package repository

import (
	"context"
	"database/sql"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type UserStatsRepository struct{}

func NewUserStatsRepository() *UserStatsRepository {
	return &UserStatsRepository{}
}

func (r *UserStatsRepository) LockForUpdate(ctx context.Context, tx *sql.Tx, userID string) (*domain.UserStats, error) {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO user_stats (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
		`, userID)
	if err != nil {
		return nil, err
	}

	var s domain.UserStats
	err = tx.QueryRowContext(ctx, `
		SELECT user_id, total_min, current_streak, longest_streak, last_studied_on, updated_at
		FROM user_stats
		WHERE user_id = $1
		FOR UPDATE
		`, userID).Scan(&s.UserID, &s.TotalMin, &s.CurrentStreak, &s.LongestStreak, &s.LastStudiedOn, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *UserStatsRepository) Save(ctx context.Context, tx *sql.Tx, s *domain.UserStats) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE user_stats
		SET total_min = $1, current_streak = $2, longest_streak = $3, last_studied_on = $4, updated_at = now()
		WHERE user_id = $5
		`, s.TotalMin, s.CurrentStreak, s.LongestStreak, s.LastStudiedOn, s.UserID)
	return err
}
