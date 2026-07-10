package service

import (
	"context"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type CategoryRepo interface {
	FindAll(ctx context.Context, userID string) ([]domain.Category, error)
	Create(ctx context.Context, c *domain.Category) error
}

type CategoryService struct {
	repo CategoryRepo
}

func NewCategoryService(repo CategoryRepo) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) List(ctx context.Context, userID string) ([]domain.Category, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *CategoryService) Create(ctx context.Context, userID, name string) (*domain.Category, error) {
	c := &domain.Category{
		Name: name,
	}
	if userID != "" {
		c.UserID = &userID
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}
