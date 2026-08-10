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

type fakeAuthServicer struct {
	signupFunc func(ctx context.Context, email, password string) (*domain.User, error)
	loginFunc  func(ctx context.Context, email, password string) (string, error)
	meFunc     func(ctx context.Context, userID string) (*domain.User, error)
}

func (f *fakeAuthServicer) Signup(ctx context.Context, email, password string) (*domain.User, error) {
	return f.signupFunc(ctx, email, password)
}

func (f *fakeAuthServicer) Login(ctx context.Context, email, password string) (string, error) {
	return f.loginFunc(ctx, email, password)
}

func (f *fakeAuthServicer) Me(ctx context.Context, userID string) (*domain.User, error) {
	return f.meFunc(ctx, userID)
}

func TestAuthHandler_Signup(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		signupFunc func(ctx context.Context, email, password string) (*domain.User, error)
		wantStatus int
	}{
		{
			name: "正常系",
			body: `{"email":"a@example.com","password":"password123"}`,
			signupFunc: func(ctx context.Context, email, password string) (*domain.User, error) {
				return &domain.User{ID: "1", Email: email}, nil
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "emailが空",
			body:       `{"email":"","password":"password123"}`,
			signupFunc: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "passwordが空",
			body:       `{"email":"a@example.com","password":""}`,
			signupFunc: nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "メール重複",
			body: `{"email":"a@example.com","password":"password123"}`,
			signupFunc: func(ctx context.Context, email, password string) (*domain.User, error) {
				return nil, domain.ErrEmailAlreadyExists
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "その他のserviceエラー",
			body: `{"email":"a@example.com","password":"password123"}`,
			signupFunc: func(ctx context.Context, email, password string) (*domain.User, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{service: &fakeAuthServicer{signupFunc: tt.signupFunc}}
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.Signup(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		loginFunc  func(ctx context.Context, email, password string) (string, error)
		wantStatus int
		wantCookie bool
	}{
		{
			name: "正常系",
			body: `{"email":"a@example.com","password":"password123"}`,
			loginFunc: func(ctx context.Context, email, password string) (string, error) {
				return "dummy-token", nil
			},
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		{
			name: "認証失敗",
			body: `{"email":"a@example.com","password":"wrong"}`,
			loginFunc: func(ctx context.Context, email, password string) (string, error) {
				return "", domain.ErrInvalidCredentials
			},
			wantStatus: http.StatusUnauthorized,
			wantCookie: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{service: &fakeAuthServicer{loginFunc: tt.loginFunc}, secureCookie: false}
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()
			h.Login(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			cookies := w.Result().Cookies()
			hasCookie := false
			for _, c := range cookies {
				if c.Name == "token" && c.Value == "dummy-token" {
					hasCookie = true
				}
			}
			if hasCookie != tt.wantCookie {
				t.Errorf("Set-Cookie present = %v, want %v", hasCookie, tt.wantCookie)
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	h := &AuthHandler{secureCookie: false}
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "token" {
			found = true
			if c.MaxAge >= 0 {
				t.Errorf("MaxAge = %d, want negative (deletion)", c.MaxAge)
			}
		}
	}
	if !found {
		t.Error("token cookie not set for deletion")
	}
}

func TestAuthHandler_Me(t *testing.T) {
	tests := []struct {
		name       string
		meFunc     func(ctx context.Context, userID string) (*domain.User, error)
		wantStatus int
	}{
		{
			name: "正常系",
			meFunc: func(ctx context.Context, userID string) (*domain.User, error) {
				return &domain.User{ID: "1", Email: "a@example.com"}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "ユーザーが存在しない",
			meFunc: func(ctx context.Context, userID string) (*domain.User, error) {
				return nil, domain.ErrUserNotFound
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "その他のserviceエラー",
			meFunc: func(ctx context.Context, userID string) (*domain.User, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{service: &fakeAuthServicer{meFunc: tt.meFunc}}
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			w := httptest.NewRecorder()
			h.Me(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
