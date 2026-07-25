package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Natti3588/go-StudyLog/backend/internal/service"
)

type StatsHandler struct {
	service *service.StatsService
}

func NewStatsHandler(s *service.StatsService) *StatsHandler {
	return &StatsHandler{service: s}
}

func (h *StatsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	summary, err := h.service.Summary(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func (h *StatsHandler) Heatmap(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "year query parameter is required and must be a number", http.StatusBadRequest)
		return
	}

	userID := UserIDFromContext(r.Context())

	totals, err := h.service.Heatmap(r.Context(), userID, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totals)
}
