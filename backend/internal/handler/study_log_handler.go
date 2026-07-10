package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/domain"
	"github.com/Natti3588/go-StudyLog/backend/internal/service"
)

type StudyLogHandler struct {
	service *service.StudyLogService
}

type StudyLogRequest struct {
	UserID      string    `json:"user_id"`
	CategoryID  string    `json:"category_id"`
	StudiedOn   time.Time `json:"studied_on"`
	DurationMin int       `json:"duration_min"`
	Memo        string    `json:"memo"`
}

func NewStudyLogHandler(s *service.StudyLogService) *StudyLogHandler {
	return &StudyLogHandler{service: s}
}

func (h *StudyLogHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	logs, err := h.service.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *StudyLogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req StudyLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.UserID == "" || req.CategoryID == "" || req.DurationMin <= 0 {
		http.Error(w, "category_id and a positive duration_min are required", http.StatusBadRequest)
		return
	}
	l, err := h.service.Create(r.Context(), req.UserID, service.StudyLogInput{
		CategoryID:  req.CategoryID,
		StudiedOn:   req.StudiedOn,
		DurationMin: req.DurationMin,
		Memo:        req.Memo,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(l)
}

func (h *StudyLogHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req StudyLogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	l, err := h.service.Update(r.Context(), id, req.UserID, service.StudyLogInput{
		CategoryID:  req.CategoryID,
		StudiedOn:   req.StudiedOn,
		DurationMin: req.DurationMin,
		Memo:        req.Memo,
	})
	if err != nil {
		if errors.Is(err, domain.ErrStudyLogNotFound) {
			http.Error(w, "study log not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(l)
}

func (h *StudyLogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.URL.Query().Get("user_id")

	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		if errors.Is(err, domain.ErrStudyLogNotFound) {
			http.Error(w, "study log not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
