package http

import (
	"encoding/json"
	"net/http"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

func (h *HTTPHandler) CreateContest(w http.ResponseWriter, r *http.Request) {
	var req CreateContestRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*domain.AuthUser)

	input, err := mapCreateContestToInput(req)
	if err != nil {
		http.Error(w, "invalid contest time format", http.StatusBadRequest)
		return
	}

	input.CreatorID = user.ID

	out, err := h.uc.CreateContest(r.Context(), input)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapCreateContestOutputToJson(out))
}

func (h *HTTPHandler) GetContestBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.Error(w, "invalid contest slug", http.StatusBadRequest)
		return
	}

	out, err := h.uc.GetContestBySlug(r.Context(), mapGetContestBySlugToInput(slug))
	if err != nil {
		http.Error(w, "contest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapGetContestBySlugOutputToJson(out))
}
