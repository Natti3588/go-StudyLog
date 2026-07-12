package service

import (
	"context"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type StudyLogRepo interface {
	FindAll(ctx context.Context, userID string) ([]domain.StudyLog, error)
	Create(ctx context.Context, l *domain.StudyLog) error
	Update(ctx context.Context, l *domain.StudyLog) error
	Delete(ctx context.Context, id, userID string) error
}

type StudyLogService struct {
	repo StudyLogRepo
}

type StudyLogInput struct {
	CategoryID  string
	StudiedOn   time.Time
	DurationMin int
	Memo        string
}

func NewStudyLogService(repo StudyLogRepo) *StudyLogService {
	return &StudyLogService{repo: repo}
}

func (s *StudyLogService) List(ctx context.Context, userID string) ([]domain.StudyLog, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *StudyLogService) Create(ctx context.Context, userID string, in StudyLogInput) (*domain.StudyLog, error) {
	l := &domain.StudyLog{
		UserID:      userID,
		CategoryID:  in.CategoryID,
		StudiedOn:   in.StudiedOn,
		DurationMin: in.DurationMin,
		Memo:        in.Memo,
	}
	if err := s.repo.Create(ctx, l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *StudyLogService) Update(ctx context.Context, id, userID string, in StudyLogInput) (*domain.StudyLog, error) {
	l := &domain.StudyLog{
		ID:          id,
		UserID:      userID,
		CategoryID:  in.CategoryID,
		StudiedOn:   in.StudiedOn,
		DurationMin: in.DurationMin,
		Memo:        in.Memo,
	}
	if err := s.repo.Update(ctx, l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *StudyLogService) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
