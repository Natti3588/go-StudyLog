package repository

import (
	"context"
	"database/sql"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll(ctx context.Context, userID string) ([]domain.Category, error) {
	var uid *string
	if userID != "" {
		uid = &userID
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, created_at
		FROM categories
		WHERE user_id IS NULL OR user_id = $1
		ORDER BY created_at
		`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []domain.Category{}
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Create(ctx context.Context, c *domain.Category) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO categories (user_id, name)
		VALUES ($1, $2)
		RETURNING id, created_at
		`, c.UserID, c.Name).Scan(&c.ID, &c.CreatedAt)
}
