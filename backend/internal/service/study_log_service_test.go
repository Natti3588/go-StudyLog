package service

import (
	"testing"
	"time"
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
