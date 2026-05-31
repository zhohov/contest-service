package http

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.uc.Login(r.Context(), mapLoginToinput(req))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapLoginToJson(out))
}

func (h *HTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequestJson
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.uc.Register(r.Context(), mapRegisterToInput(req))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapRegisterToJson(out))
}
