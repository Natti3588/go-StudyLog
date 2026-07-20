package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type StudyLogRepository struct {
	db *sql.DB
}

func NewStudyLogRepository(db *sql.DB) *StudyLogRepository {
	return &StudyLogRepository{db: db}
}

func (r *StudyLogRepository) FindAll(ctx context.Context, userID string) ([]domain.StudyLog, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, category_id, studied_on, duration_min, memo, created_at, updated_at
		FROM study_logs
		WHERE user_id = $1
		ORDER BY studied_on DESC
		`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []domain.StudyLog{}
	for rows.Next() {
		var l domain.StudyLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.CategoryID, &l.StudiedOn, &l.DurationMin, &l.Memo, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *StudyLogRepository) Create(ctx context.Context, l *domain.StudyLog) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO study_logs (user_id, category_id, studied_on, duration_min, memo)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
		`, l.UserID, l.CategoryID, l.StudiedOn, l.DurationMin, l.Memo).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
}

func (r *StudyLogRepository) Update(ctx context.Context, l *domain.StudyLog) error {
	err := r.db.QueryRowContext(ctx, `
		UPDATE study_logs
		SET category_id = $1, studied_on = $2, duration_min = $3, memo = $4, updated_at = now()
		WHERE id = $5 AND user_id = $6
		RETURNING created_at, updated_at
		`, l.CategoryID, l.StudiedOn, l.DurationMin, l.Memo, l.ID, l.UserID).Scan(&l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrStudyLogNotFound
		}
		return err
	}
	return nil
}

func (r *StudyLogRepository) Delete(ctx context.Context, id, userID string) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM study_logs WHERE id = $1 AND user_id = $2
		`, id, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrStudyLogNotFound
	}
	return nil
}

func (r *StudyLogRepository) CreateTx(ctx context.Context, tx *sql.Tx, l *domain.StudyLog) error {
	return tx.QueryRowContext(ctx, `
        INSERT INTO study_logs (user_id, category_id, studied_on, duration_min, memo)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
		`, l.UserID, l.CategoryID, l.StudiedOn, l.DurationMin, l.Memo).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
}

func (r *StudyLogRepository) FindAllByUserTx(ctx context.Context, tx *sql.Tx, userID string) ([]domain.StudyLog, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, user_id, category_id, studied_on, duration_min, memo, created_at, updated_at
		FROM study_logs
		WHERE user_id = $1
		ORDER BY studied_on
		`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []domain.StudyLog{}
	for rows.Next() {
		var l domain.StudyLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.CategoryID, &l.StudiedOn, &l.DurationMin, &l.Memo, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

func (r *StudyLogRepository) UpdateTx(ctx context.Context, tx *sql.Tx, l *domain.StudyLog) error {
	err := tx.QueryRowContext(ctx, `
		UPDATE study_logs
		SET category_id = $1, studied_on = $2, duration_min = $3, memo = $4, updated_at = now()
		WHERE id = $5 AND user_id = $6
		RETURNING created_at, updated_at
	`, l.CategoryID, l.StudiedOn, l.DurationMin, l.Memo, l.ID, l.UserID).Scan(&l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrStudyLogNotFound
		}
		return err
	}
	return nil
}

func (r *StudyLogRepository) DeleteTx(ctx context.Context, tx *sql.Tx, id, userID string) error {
	res, err := tx.ExecContext(ctx, `DELETE FROM study_logs WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrStudyLogNotFound
	}
	return nil
}
