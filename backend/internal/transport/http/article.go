package http

import (
	"encoding/json"
	"net/http"

	"github.com/zhohov/contest-service/backend/internal/domain"
)

func (h *HTTPHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req CreateArticleRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*domain.AuthUser)

	input := mapCreateArticleToInput(req)
	input.AuthorID = user.ID

	out, err := h.uc.CreateArticle(r.Context(), input)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapCreateArticleOutputToJson(out))
}

func (h *HTTPHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	out, err := h.uc.GetArticles(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mapGetArticlesOutputToJson(out))
}

// func (h *HTTPHandler) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {

// }
