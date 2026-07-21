package service

import (
	"context"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type WeeklyGoalRepo interface {
	Upsert(ctx context.Context, g *domain.WeeklyGoal) error
	FindByUserAndWeek(ctx context.Context, userID string, week time.Time) (*domain.WeeklyGoal, error)
}

type WeeklyGoalService struct {
	repo WeeklyGoalRepo
}

func NewWeeklyGoalService(repo WeeklyGoalRepo) *WeeklyGoalService {
	return &WeeklyGoalService{repo: repo}
}

func (s *WeeklyGoalService) SetWeekly(ctx context.Context, userID string, weekStart time.Time, targetMin int) (*domain.WeeklyGoal, error) {
	g := &domain.WeeklyGoal{
		UserID:    userID,
		WeekStart: weekStart,
		TargetMin: targetMin,
	}
	if err := s.repo.Upsert(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}
