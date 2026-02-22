package handler

import (
	"net/http"

	"baby-care/internal/model"
	"baby-care/internal/store"
)

type feedingRequest struct {
	FeedType   string `json:"feed_type"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	QuantityML *int   `json:"quantity_ml"`
	Notes      string `json:"notes"`
}

func (h *Handler) ListFeeding(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	date := r.URL.Query().Get("date")
	logs, err := h.Store.GetFeedingLogs(childID, date)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if logs == nil {
		logs = []*model.FeedingLog{}
	}
	h.JSON(w, http.StatusOK, logs)
}

func (h *Handler) CreateFeeding(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	var req feedingRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.FeedType == "" {
		h.Error(w, http.StatusBadRequest, "feed_type is required")
		return
	}
	log, stopped, err := h.Store.CreateFeeding(childID, req.FeedType, req.StartTime, req.Notes, req.QuantityML)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	type createFeedingResponse struct {
		*model.FeedingLog
		StoppedSleep *store.StoppedSleep `json:"stopped_sleep,omitempty"`
	}
	h.JSON(w, http.StatusCreated, createFeedingResponse{FeedingLog: log, StoppedSleep: stopped})
}

func (h *Handler) GetActiveFeeding(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	log, err := h.Store.GetActiveFeeding(childID)
	if err != nil {
		if h.IsNotFound(err) {
			h.JSON(w, http.StatusOK, nil)
			return
		}
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, log)
}

func (h *Handler) UpdateFeeding(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	var req feedingRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	log, err := h.Store.UpdateFeeding(id, req.FeedType, req.StartTime, req.EndTime, req.Notes, req.QuantityML)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, log)
}

func (h *Handler) DeleteFeeding(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	if err := h.Store.DeleteFeeding(id); err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
