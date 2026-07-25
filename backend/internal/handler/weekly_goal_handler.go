package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Natti3588/go-StudyLog/backend/internal/service"
)

type WeeklyGoalHandler struct {
	service *service.WeeklyGoalService
}

type weeklyGoalRequest struct {
	WeekStart time.Time `json:"week_start"`
	TargetMin int       `json:"target_min"`
}

func NewWeeklyGoalHandler(s *service.WeeklyGoalService) *WeeklyGoalHandler {
	return &WeeklyGoalHandler{service: s}
}

func (h *WeeklyGoalHandler) SetWeekly(w http.ResponseWriter, r *http.Request) {
	var req weeklyGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.TargetMin <= 0 {
		http.Error(w, "target_min must be positive", http.StatusBadRequest)
		return
	}

	userID := UserIDFromContext(r.Context())
	g, err := h.service.SetWeekly(r.Context(), userID, req.WeekStart, req.TargetMin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}
