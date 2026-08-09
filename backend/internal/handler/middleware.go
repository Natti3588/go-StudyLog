package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/Natti3588/go-StudyLog/backend/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func RequireAuth(jwtSecret []byte) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractToken(r)
			if tokenString == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims := &service.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next(w, r.WithContext(ctx))
		}
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func UserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(userIDContextKey).(string)
	return userID
}

func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
