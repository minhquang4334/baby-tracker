package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"baby-care/internal/store"
)

type Handler struct {
	Store *store.Store
}

func (h *Handler) JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) Error(w http.ResponseWriter, status int, msg string) {
	h.JSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) Decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *Handler) IsNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}
