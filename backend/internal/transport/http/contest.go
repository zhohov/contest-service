package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/zhohov/contest-service/backend/internal/domain"
	"github.com/zhohov/contest-service/backend/internal/usecase"
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

func (h *HTTPHandler) GetContestsByUserID(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(*domain.AuthUser)

	out, err := h.uc.GetContestsByUserID(r.Context(), mapGetContestsByUserIdToInput(user.ID))
	if err != nil {
		http.Error(w, "contests not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapGetContestsByUserIdOutputToJson(out))
}

func (h *HTTPHandler) UpdateContest(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "invalid contest id", http.StatusBadRequest)
		return
	}

	var req UpdateContestRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*domain.AuthUser)

	input, err := mapUpdateContestToInput(id, user.ID, req)
	if err != nil {
		http.Error(w, "invalid contest time format", http.StatusBadRequest)
		return
	}

	out, err := h.uc.UpdateContest(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrContestNotFound):
			http.Error(w, "contest not found", http.StatusNotFound)
		case errors.Is(err, usecase.ErrContestForbidden):
			http.Error(w, "forbidden", http.StatusForbidden)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapUpdateContestOutputToJson(out))
}
