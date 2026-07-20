package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type StudyLogRepo interface {
	FindAll(ctx context.Context, userID string) ([]domain.StudyLog, error)
	Create(ctx context.Context, l *domain.StudyLog) error
	Update(ctx context.Context, l *domain.StudyLog) error
	Delete(ctx context.Context, id, userID string) error
	CreateTx(ctx context.Context, tx *sql.Tx, l *domain.StudyLog) error
	FindAllByUserTx(ctx context.Context, tx *sql.Tx, userID string) ([]domain.StudyLog, error)
	UpdateTx(ctx context.Context, tx *sql.Tx, l *domain.StudyLog) error
	DeleteTx(ctx context.Context, tx *sql.Tx, id, userID string) error
}

type UserStatsRepo interface {
	LockForUpdate(ctx context.Context, tx *sql.Tx, userID string) (*domain.UserStats, error)
	Save(ctx context.Context, tx *sql.Tx, s *domain.UserStats) error
}

type TxBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type StudyLogService struct {
	repo      StudyLogRepo
	statsRepo UserStatsRepo
	db        TxBeginner
}

type StudyLogInput struct {
	CategoryID  string
	StudiedOn   time.Time
	DurationMin int
	Memo        string
}

func NewStudyLogService(repo StudyLogRepo, statsRepo UserStatsRepo, db TxBeginner) *StudyLogService {
	return &StudyLogService{repo: repo, statsRepo: statsRepo, db: db}
}

func (s *StudyLogService) List(ctx context.Context, userID string) ([]domain.StudyLog, error) {
	return s.repo.FindAll(ctx, userID)
}

func (s *StudyLogService) Create(ctx context.Context, userID string, in StudyLogInput) (*domain.StudyLog, error) {
	if in.StudiedOn.After(time.Now()) {
		return nil, domain.ErrInvalidStudiedOn
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stats, err := s.statsRepo.LockForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	l := &domain.StudyLog{
		UserID:      userID,
		CategoryID:  in.CategoryID,
		StudiedOn:   in.StudiedOn,
		DurationMin: in.DurationMin,
		Memo:        in.Memo,
	}
	if err := s.repo.CreateTx(ctx, tx, l); err != nil {
		return nil, err
	}

	if stats.LastStudiedOn != nil && in.StudiedOn.Before(*stats.LastStudiedOn) {
		logs, err := s.repo.FindAllByUserTx(ctx, tx, userID)
		if err != nil {
			return nil, err
		}
		recomputeStats(stats, logs)
	} else {
		applyStreakForNewLog(stats, in.StudiedOn)
		stats.TotalMin += in.DurationMin
	}

	if err := s.statsRepo.Save(ctx, tx, stats); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *StudyLogService) Update(ctx context.Context, id, userID string, in StudyLogInput) (*domain.StudyLog, error) {
	if in.StudiedOn.After(time.Now()) {
		return nil, domain.ErrInvalidStudiedOn
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	l := &domain.StudyLog{
		ID:          id,
		UserID:      userID,
		CategoryID:  in.CategoryID,
		StudiedOn:   in.StudiedOn,
		DurationMin: in.DurationMin,
		Memo:        in.Memo,
	}
	if err := s.repo.UpdateTx(ctx, tx, l); err != nil {
		return nil, err
	}

	stats, err := s.statsRepo.LockForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	logs, err := s.repo.FindAllByUserTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	recomputeStats(stats, logs)

	if err := s.statsRepo.Save(ctx, tx, stats); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return l, nil
}

func (s *StudyLogService) Delete(ctx context.Context, id, userID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.DeleteTx(ctx, tx, id, userID); err != nil {
		return err
	}

	stats, err := s.statsRepo.LockForUpdate(ctx, tx, userID)
	if err != nil {
		return err
	}
	logs, err := s.repo.FindAllByUserTx(ctx, tx, userID)
	if err != nil {
		return err
	}
	recomputeStats(stats, logs)

	if err := s.statsRepo.Save(ctx, tx, stats); err != nil {
		return err
	}

	return tx.Commit()
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
