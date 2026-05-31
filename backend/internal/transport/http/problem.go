package http

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var req CreateProblemRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.uc.CreateProblem(r.Context(), mapCreateProblemToInput(req))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapCreateProblemOutputToJson(out))
}

func (h *HTTPHandler) GetProblems(w http.ResponseWriter, r *http.Request) {
	out, err := h.uc.GetProblems(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapGetProblemsOutputToJson(out))
}

func (h *HTTPHandler) GetProblemBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	out, err := h.uc.GetProblemBySlug(r.Context(), mapGetProblemBySlugToInput(slug))
	if err != nil {
		http.Error(w, "problem not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapGetProblemBySlugOutputToJson(out))
}
