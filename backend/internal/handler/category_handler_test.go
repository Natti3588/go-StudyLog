package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
)

type fakeCategoryServicer struct {
	listFunc   func(ctx context.Context, userID string) ([]domain.Category, error)
	createFunc func(ctx context.Context, userID, name string) (*domain.Category, error)
}

func (f *fakeCategoryServicer) List(ctx context.Context, userID string) ([]domain.Category, error) {
	return f.listFunc(ctx, userID)
}

func (f *fakeCategoryServicer) Create(ctx context.Context, userID, name string) (*domain.Category, error) {
	return f.createFunc(ctx, userID, name)
}

func TestCategoryHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		listFunc   func(ctx context.Context, userID string) ([]domain.Category, error)
		wantStatus int
	}{
		{
			name: "正常系",
			listFunc: func(ctx context.Context, userID string) ([]domain.Category, error) {
				return []domain.Category{{ID: "1", Name: "リーディング"}}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "serviceがエラーを返す",
			listFunc: func(ctx context.Context, userID string) ([]domain.Category, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CategoryHandler{service: &fakeCategoryServicer{listFunc: tt.listFunc}}

			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			w := httptest.NewRecorder()

			h.List(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestCategoryHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		createFunc func(ctx context.Context, userID, name string) (*domain.Category, error)
		wantStatus int
	}{
		{
			name: "正常系",
			body: `{"name":"リーディング"}`,
			createFunc: func(ctx context.Context, userID, name string) (*domain.Category, error) {
				return &domain.Category{ID: "1", Name: name}, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "nameが空文字列",
			body:       `{"name":""}`,
			createFunc: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "serviceがエラーを返す",
			body: `{"name":"リーディング"}`,
			createFunc: func(ctx context.Context, userID, name string) (*domain.Category, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CategoryHandler{service: &fakeCategoryServicer{createFunc: tt.createFunc}}

			req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
