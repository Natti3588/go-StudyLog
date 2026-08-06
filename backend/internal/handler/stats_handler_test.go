package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type fakeStatsServicer struct {
	summaryFunc func(ctx context.Context, userID string) (*domain.StatsSummary, error)
	heatmapFunc func(ctx context.Context, userID string, year int) ([]domain.DailyTotal, error)
}

func (f *fakeStatsServicer) Summary(ctx context.Context, userID string) (*domain.StatsSummary, error) {
	return f.summaryFunc(ctx, userID)
}
func (f *fakeStatsServicer) Heatmap(ctx context.Context, userID string, year int) ([]domain.DailyTotal, error) {
	return f.heatmapFunc(ctx, userID, year)
}

func TestStatsHandler_Summary(t *testing.T) {
	tests := []struct {
		name        string
		summaryFunc func(ctx context.Context, userID string) (*domain.StatsSummary, error)
		wantStatus  int
	}{
		{
			name: "正常系",
			summaryFunc: func(ctx context.Context, userID string) (*domain.StatsSummary, error) {
				return &domain.StatsSummary{TotalMin: 100}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "serviceがエラーを返す",
			summaryFunc: func(ctx context.Context, userID string) (*domain.StatsSummary, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatsHandler{service: &fakeStatsServicer{summaryFunc: tt.summaryFunc}}
			req := httptest.NewRequest(http.MethodGet, "/stats/summary", nil)
			w := httptest.NewRecorder()
			h.Summary(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestStatsHandler_Heatmap(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		heatmapFunc func(ctx context.Context, userID string, year int) ([]domain.DailyTotal, error)
		wantStatus  int
	}{
		{
			name: "正常系",
			url:  "/stats/heatmap?year=2026",
			heatmapFunc: func(ctx context.Context, userID string, year int) ([]domain.DailyTotal, error) {
				return []domain.DailyTotal{{TotalMin: 30}}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:        "yearパラメータが無い",
			url:         "/stats/heatmap",
			heatmapFunc: nil,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name:        "yearパラメータが数値でない",
			url:         "/stats/heatmap?year=abc",
			heatmapFunc: nil,
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "serviceがエラーを返す",
			url:  "/stats/heatmap?year=2026",
			heatmapFunc: func(ctx context.Context, userID string, year int) ([]domain.DailyTotal, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatsHandler{service: &fakeStatsServicer{heatmapFunc: tt.heatmapFunc}}
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			h.Heatmap(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
