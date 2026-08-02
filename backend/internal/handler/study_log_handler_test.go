package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
	"github.com/Natti3588/go-StudyLog/backend/internal/service"
)

type fakeStudyLogServicer struct {
	listFunc   func(ctx context.Context, userID string) ([]domain.StudyLog, error)
	createFunc func(ctx context.Context, userID string, in service.StudyLogInput) (*domain.StudyLog, error)
	updateFunc func(ctx context.Context, id, userID string, in service.StudyLogInput) (*domain.StudyLog, error)
	deleteFunc func(ctx context.Context, id, userID string) error
}

func (f *fakeStudyLogServicer) List(ctx context.Context, userID string) ([]domain.StudyLog, error) {
	return f.listFunc(ctx, userID)
}

func (f *fakeStudyLogServicer) Create(ctx context.Context, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
	return f.createFunc(ctx, userID, in)
}

func (f *fakeStudyLogServicer) Update(ctx context.Context, id, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
	return f.updateFunc(ctx, id, userID, in)
}

func (f *fakeStudyLogServicer) Delete(ctx context.Context, id, userID string) error {
	return f.deleteFunc(ctx, id, userID)
}

func TestStudyLogHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		listFunc   func(ctx context.Context, userID string) ([]domain.StudyLog, error)
		wantStatus int
	}{
		{
			name: "正常系",
			listFunc: func(ctx context.Context, userID string) ([]domain.StudyLog, error) {
				return []domain.StudyLog{{ID: "1"}}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "serviceがエラーを返す",
			listFunc: func(ctx context.Context, userID string) ([]domain.StudyLog, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudyLogHandler{service: &fakeStudyLogServicer{listFunc: tt.listFunc}}
			req := httptest.NewRequest(http.MethodGet, "/logs", nil)
			w := httptest.NewRecorder()
			h.List(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestStudyLogHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		createFunc func(ctx context.Context, userID string, in service.StudyLogInput) (*domain.StudyLog, error)
		wantStatus int
	}{
		{
			name: "正常系",
			body: `{"category_id":"c1","studied_on":"2026-07-10T00:00:00Z","duration_min":30}`,
			createFunc: func(ctx context.Context, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return &domain.StudyLog{ID: "1"}, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "category_idが空",
			body:       `{"category_id":"","duration_min":30}`,
			createFunc: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duration_minが0以下",
			body:       `{"category_id":"c1","duration_min":0}`,
			createFunc: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "未来日",
			body: `{"category_id":"c1","studied_on":"2099-01-01T00:00:00Z","duration_min":30}`,
			createFunc: func(ctx context.Context, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return nil, domain.ErrInvalidStudiedOn
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "category_idが存在しない",
			body: `{"category_id":"c1","duration_min":30}`,
			createFunc: func(ctx context.Context, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return nil, domain.ErrCategoryNotFound
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "その他のserviceエラー",
			body: `{"category_id":"c1","duration_min":30}`,
			createFunc: func(ctx context.Context, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudyLogHandler{
				service: &fakeStudyLogServicer{
					createFunc: tt.createFunc,
				},
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/logs",
				bytes.NewBufferString(tt.body),
			)
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestStudyLogHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		updateFunc func(ctx context.Context, id, userID string, in service.StudyLogInput) (*domain.StudyLog, error)
		wantStatus int
	}{
		{
			name: "正常系",
			body: `{"category_id":"c1","duration_min":30}`,
			updateFunc: func(ctx context.Context, id, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return &domain.StudyLog{ID: id}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "対象が存在しない",
			body: `{"category_id":"c1","duration_min":30}`,
			updateFunc: func(ctx context.Context, id, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return nil, domain.ErrStudyLogNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "未来日",
			body: `{"category_id":"c1","duration_min":30}`,
			updateFunc: func(ctx context.Context, id, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return nil, domain.ErrInvalidStudiedOn
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "category_idが存在しない",
			body: `{"category_id":"c1","duration_min":30}`,
			updateFunc: func(ctx context.Context, id, userID string, in service.StudyLogInput) (*domain.StudyLog, error) {
				return nil, domain.ErrCategoryNotFound
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudyLogHandler{service: &fakeStudyLogServicer{updateFunc: tt.updateFunc}}
			req := httptest.NewRequest(http.MethodPut, "/logs/1", bytes.NewBufferString(tt.body))
			req.SetPathValue("id", "1")
			w := httptest.NewRecorder()
			h.Update(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestStudyLogHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		deleteFunc func(ctx context.Context, id, userID string) error
		wantStatus int
	}{
		{
			name: "正常系",
			deleteFunc: func(ctx context.Context, id, userID string) error {
				return nil
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "対象が存在しない",
			deleteFunc: func(ctx context.Context, id, userID string) error {
				return domain.ErrStudyLogNotFound
			},
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudyLogHandler{service: &fakeStudyLogServicer{deleteFunc: tt.deleteFunc}}
			req := httptest.NewRequest(http.MethodDelete, "/logs/1", nil)
			req.SetPathValue("id", "1")
			w := httptest.NewRecorder()
			h.Delete(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
