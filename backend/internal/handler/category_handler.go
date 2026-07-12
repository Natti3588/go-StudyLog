package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Natti3588/go-StudyLog/backend/internal/service"
)

type CategoryHandler struct {
	service *service.CategoryService
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())

	categories, err := h.service.List(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	userID := UserIDFromContext(r.Context())
	c, err := h.service.Create(r.Context(), userID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}
