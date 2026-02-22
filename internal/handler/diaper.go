package handler

import (
	"net/http"

	"baby-care/internal/model"
)

type diaperRequest struct {
	DiaperType string `json:"diaper_type"`
	ChangedAt  string `json:"changed_at"`
	Notes      string `json:"notes"`
}

func (h *Handler) ListDiaper(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	date := r.URL.Query().Get("date")
	logs, err := h.Store.GetDiaperLogs(childID, date)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if logs == nil {
		logs = []*model.DiaperLog{}
	}
	h.JSON(w, http.StatusOK, logs)
}

func (h *Handler) CreateDiaper(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	var req diaperRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.DiaperType == "" {
		h.Error(w, http.StatusBadRequest, "diaper_type is required")
		return
	}
	log, err := h.Store.CreateDiaper(childID, req.DiaperType, req.ChangedAt, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusCreated, log)
}

func (h *Handler) UpdateDiaper(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	var req diaperRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	log, err := h.Store.UpdateDiaper(id, req.DiaperType, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, log)
}

func (h *Handler) DeleteDiaper(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	if err := h.Store.DeleteDiaper(id); err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
