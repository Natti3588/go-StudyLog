package service

import (
	"context"
	"errors"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
}

type AuthService struct {
	repo      UserRepo
	jwtSecret []byte
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAuthService(repo UserRepo, jwtSecret []byte) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (*domain.User, error) {
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrEmailAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	claims := Claims{
		UserID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
