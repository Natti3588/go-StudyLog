package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, is_admin, created_at
		FROM users
		WHERE email = $1
		`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, is_admin, created_at
		`, u.Email, u.PasswordHash).Scan(&u.ID, &u.IsAdmin, &u.CreatedAt)
}
