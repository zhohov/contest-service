package http

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var req CreateArticleRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.uc.CreateArticle(mapCreateArticleToInput(req))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapCreateArticleOutputToJson(out))
}

func (h *HTTPHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	out, err := h.uc.GetArticles()
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
