package handler

import (
	"net/http"

	"baby-care/internal/model"
)

type growthRequest struct {
	MeasuredOn          string `json:"measured_on"`
	WeightGrams         *int   `json:"weight_grams"`
	LengthMM            *int   `json:"length_mm"`
	HeadCircumferenceMM *int   `json:"head_circumference_mm"`
	Notes               string `json:"notes"`
}

func (h *Handler) ListGrowth(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	logs, err := h.Store.GetGrowthLogs(childID)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if logs == nil {
		logs = []*model.GrowthLog{}
	}
	h.JSON(w, http.StatusOK, logs)
}

func (h *Handler) CreateGrowth(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	var req growthRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	log, err := h.Store.CreateGrowth(childID, req.MeasuredOn, req.WeightGrams, req.LengthMM, req.HeadCircumferenceMM, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusCreated, log)
}

func (h *Handler) UpdateGrowth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	var req growthRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	log, err := h.Store.UpdateGrowth(id, req.MeasuredOn, req.WeightGrams, req.LengthMM, req.HeadCircumferenceMM, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, log)
}

func (h *Handler) DeleteGrowth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	if err := h.Store.DeleteGrowth(id); err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
