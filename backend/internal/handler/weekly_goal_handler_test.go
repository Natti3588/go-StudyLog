package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type fakeWeeklyGoalServicer struct {
	setWeeklyFunc func(ctx context.Context, userID string, weekStart time.Time, targetMin int) (*domain.WeeklyGoal, error)
}

func (f *fakeWeeklyGoalServicer) SetWeekly(ctx context.Context, userID string, weekStart time.Time, targetMin int) (*domain.WeeklyGoal, error) {
	return f.setWeeklyFunc(ctx, userID, weekStart, targetMin)
}

func TestWeeklyGoalHandler_SetWeekly(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		setWeeklyFunc func(ctx context.Context, userID string, weekStart time.Time, targetMin int) (*domain.WeeklyGoal, error)
		wantStatus    int
	}{
		{
			name: "正常系",
			body: `{"week_start":"2026-07-20T00:00:00Z","target_min":200}`,
			setWeeklyFunc: func(ctx context.Context, userID string, weekStart time.Time, targetMin int) (*domain.WeeklyGoal, error) {
				return &domain.WeeklyGoal{UserID: userID, WeekStart: weekStart, TargetMin: targetMin}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:          "target_minが0以下",
			body:          `{"week_start":"2026-07-20T00:00:00Z","target_min":0}`,
			setWeeklyFunc: nil,
			wantStatus:    http.StatusBadRequest,
		},
		{
			name: "serviceがエラーを返す",
			body: `{"week_start":"2026-07-20T00:00:00Z","target_min":200}`,
			setWeeklyFunc: func(ctx context.Context, userID string, weekStart time.Time, targetMin int) (*domain.WeeklyGoal, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &WeeklyGoalHandler{service: &fakeWeeklyGoalServicer{setWeeklyFunc: tt.setWeeklyFunc}}
			req := httptest.NewRequest(http.MethodPost, "/goals/weekly", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			h.SetWeekly(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
