package service

import (
	"context"
	"errors"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type StatsLogRepo interface {
	FindDailyTotals(ctx context.Context, userID string, year int) ([]domain.DailyTotal, error)
	SumDurationInRange(ctx context.Context, userID string, from, to time.Time) (int, error)
}

type StatsUserStatsRepo interface {
	FindByUserID(ctx context.Context, userID string) (*domain.UserStats, error)
}

type StatsGoalRepo interface {
	FindByUserAndWeek(ctx context.Context, userID string, weekStart time.Time) (*domain.WeeklyGoal, error)
}

type StatsService struct {
	logRepo   StatsLogRepo
	statsRepo StatsUserStatsRepo
	goalRepo  StatsGoalRepo
}

func NewStatsService(logRepo StatsLogRepo, statsRepo StatsUserStatsRepo, goalRepo StatsGoalRepo) *StatsService {
	return &StatsService{logRepo: logRepo, statsRepo: statsRepo, goalRepo: goalRepo}
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	y, m, d := t.AddDate(0, 0, -(weekday - 1)).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (s *StatsService) Summary(ctx context.Context, userID string) (*domain.StatsSummary, error) {
	stats, err := s.statsRepo.FindByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, domain.ErrUserStatsNotFound) {
			return nil, err
		}
		stats = &domain.UserStats{}
	}

	weekStart := startOfWeek(time.Now())
	weekEnd := weekStart.AddDate(0, 0, 7)

	weeklyActual, err := s.logRepo.SumDurationInRange(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	goal, err := s.goalRepo.FindByUserAndWeek(ctx, userID, weekStart)
	weeklyTarget := 0
	if err != nil {
		if !errors.Is(err, domain.ErrWeeklyGoalNotFound) {
			return nil, err
		}
	} else {
		weeklyTarget = goal.TargetMin
	}

	return &domain.StatsSummary{
		TotalMin:        stats.TotalMin,
		CurrentStreak:   stats.CurrentStreak,
		LongestStreak:   stats.LongestStreak,
		WeeklyTargetMin: weeklyTarget,
		WeeklyActualMin: weeklyActual,
	}, nil
}

func (s *StatsService) Heatmap(ctx context.Context, userID string, year int) ([]domain.DailyTotal, error) {
	return s.logRepo.FindDailyTotals(ctx, userID, year)
}
