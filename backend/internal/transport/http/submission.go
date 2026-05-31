package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

type Event struct {
	Message string `json:"message"`
}

func (h *HTTPHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var req CreateSubmissionRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*domain.AuthUser)
	log.Println(user)

	input := mapCreateSubmissionToInput(req)
	input.UserID = user.ID

	out, err := h.uc.CreateSubmission(r.Context(), input)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapCreateSubmissionOutputToJson(out))
}

func (h *HTTPHandler) GetSubmissionByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid submission id", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey)
	log.Println(user)

	out, err := h.uc.GetSubmissionByID(r.Context(), mapGetSubmissionByIdToInput(id))
	if err != nil {
		http.Error(w, "submission not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapGetSubmissionByIdOutputToJson(out))
}

// func (h *HTTPHandler) GetSubmissionByUserID(w http.ResponseWriter, r *http.Request) {

// }

// func (h *HTTPHandler) GetSubmissionByContestID(w http.ResponseWriter, r *http.Request) {

// }
