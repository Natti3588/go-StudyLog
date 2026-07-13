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

func daysBetween(a, b time.Time) int {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	da := time.Date(ay, am, ad, 0, 0, 0, 0, time.UTC)
	db := time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	return int(db.Sub(da).Hours() / 24)
}

func applyStreakForNewLog(stats *domain.UserStats, studiedOn time.Time) {
	if stats.LastStudiedOn == nil {
		stats.CurrentStreak = 1
	} else {
		switch daysBetween(*stats.LastStudiedOn, studiedOn) {
		case 0:
			// 同日の2件目以降は変化なし
		case 1:
			stats.CurrentStreak++
		default:
			stats.CurrentStreak = 1
		}
	}
	if stats.CurrentStreak > stats.LongestStreak {
		stats.LongestStreak = stats.CurrentStreak
	}
	stats.LastStudiedOn = &studiedOn
}

func recomputeStats(stats *domain.UserStats, logs []domain.StudyLog) {
	totalMin := 0
	var prevDay *time.Time
	streak := 0
	longest := 0

	for _, l := range logs {
		totalMin += l.DurationMin

		if prevDay == nil {
			streak = 1
		} else {
			switch daysBetween(*prevDay, l.StudiedOn) {
			case 0:
				// 同日は変化なし
			case 1:
				streak++
			default:
				streak = 1
			}
		}

		if streak > longest {
			longest = streak
		}

		day := l.StudiedOn
		prevDay = &day
	}

	stats.TotalMin = totalMin
	stats.CurrentStreak = streak
	stats.LongestStreak = longest
	stats.LastStudiedOn = prevDay
}
