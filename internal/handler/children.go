package handler

import (
	"net/http"
)

type childRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	PhotoURL    string `json:"photo_url"`
	Notes       string `json:"notes"`
}

func (h *Handler) GetChild(w http.ResponseWriter, r *http.Request) {
	child, err := h.Store.GetChild()
	if err != nil {
		if h.IsNotFound(err) {
			h.JSON(w, http.StatusNotFound, map[string]string{"message": "no child yet"})
			return
		}
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, child)
}

func (h *Handler) CreateChild(w http.ResponseWriter, r *http.Request) {
	var req childRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Name == "" || req.DateOfBirth == "" || req.Gender == "" {
		h.Error(w, http.StatusBadRequest, "name, date_of_birth, and gender are required")
		return
	}
	child, err := h.Store.CreateChild(req.Name, req.DateOfBirth, req.Gender, req.PhotoURL, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusCreated, child)
}

func (h *Handler) UpdateChild(w http.ResponseWriter, r *http.Request) {
	existing, err := h.Store.GetChild()
	if err != nil {
		if h.IsNotFound(err) {
			h.Error(w, http.StatusNotFound, "no child profile found")
			return
		}
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	var req childRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	child, err := h.Store.UpdateChild(existing.ID, req.Name, req.DateOfBirth, req.Gender, req.PhotoURL, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, child)
}
