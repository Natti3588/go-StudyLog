package service

import (
	"testing"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

func TestDaysBetween(t *testing.T) {
	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want int
	}{
		{
			name: "同じ日",
			a:    time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		{
			name: "連続した翌日",
			a:    time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC),
			want: 1,
		},
		{
			name: "3日空いている",
			a:    time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2026, 7, 13, 0, 0, 0, 0, time.UTC),
			want: 3,
		},
		{
			name: "時刻が違っても日付だけで比較する",
			a:    time.Date(2026, 7, 10, 23, 59, 0, 0, time.UTC),
			b:    time.Date(2026, 7, 11, 0, 1, 0, 0, time.UTC),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := daysBetween(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("daysBetween(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func day(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func ptr(t time.Time) *time.Time {
	return &t
}

func TestApplyStreakForNewLog(t *testing.T) {
	tests := []struct {
		name              string
		initial           domain.UserStats
		studiedOn         time.Time
		wantCurrentStreak int
		wantLongestStreak int
		wantLastStudiedOn time.Time
	}{
		{
			name:              "初回投稿",
			initial:           domain.UserStats{CurrentStreak: 0, LongestStreak: 0, LastStudiedOn: nil},
			studiedOn:         day(2026, 7, 10),
			wantCurrentStreak: 1,
			wantLongestStreak: 1,
			wantLastStudiedOn: day(2026, 7, 10),
		},
		{
			name: "同日の2件目は変化しない",
			initial: domain.UserStats{
				CurrentStreak: 3, LongestStreak: 5,
				LastStudiedOn: ptr(day(2026, 7, 10)),
			},
			studiedOn:         day(2026, 7, 10),
			wantCurrentStreak: 3,
			wantLongestStreak: 5,
			wantLastStudiedOn: day(2026, 7, 10),
		},
		{
			name: "連続した翌日は+1",
			initial: domain.UserStats{
				CurrentStreak: 3, LongestStreak: 3,
				LastStudiedOn: ptr(day(2026, 7, 10)),
			},
			studiedOn:         day(2026, 7, 11),
			wantCurrentStreak: 4,
			wantLongestStreak: 4,
			wantLastStudiedOn: day(2026, 7, 11),
		},
		{
			name: "空きがあった場合は1にリセット、longestは維持",
			initial: domain.UserStats{
				CurrentStreak: 5, LongestStreak: 10,
				LastStudiedOn: ptr(day(2026, 7, 10)),
			},
			studiedOn:         day(2026, 7, 13),
			wantCurrentStreak: 1,
			wantLongestStreak: 10,
			wantLastStudiedOn: day(2026, 7, 13),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := tt.initial
			applyStreakForNewLog(&stats, tt.studiedOn)

			if stats.CurrentStreak != tt.wantCurrentStreak {
				t.Errorf("CurrentStreak = %d, want %d", stats.CurrentStreak, tt.wantCurrentStreak)
			}
			if stats.LongestStreak != tt.wantLongestStreak {
				t.Errorf("LongestStreak = %d, want %d", stats.LongestStreak, tt.wantLongestStreak)
			}
			if stats.LastStudiedOn == nil || !stats.LastStudiedOn.Equal(tt.wantLastStudiedOn) {
				t.Errorf("LastStudiedOn = %v, want %v", stats.LastStudiedOn, tt.wantLastStudiedOn)
			}
		})
	}
}
