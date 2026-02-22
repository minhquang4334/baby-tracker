package handler

import (
	"net/http"

	"baby-care/internal/model"
	"baby-care/internal/store"
)

type sleepRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Notes     string `json:"notes"`
}

func (h *Handler) requireChild(w http.ResponseWriter) (string, bool) {
	child, err := h.Store.GetChild()
	if err != nil {
		if h.IsNotFound(err) {
			h.Error(w, http.StatusBadRequest, "no child profile found")
		} else {
			h.Error(w, http.StatusInternalServerError, err.Error())
		}
		return "", false
	}
	return child.ID, true
}

func (h *Handler) ListSleep(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	date := r.URL.Query().Get("date")
	logs, err := h.Store.GetSleepLogs(childID, date)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if logs == nil {
		logs = []*model.SleepLog{}
	}
	h.JSON(w, http.StatusOK, logs)
}

func (h *Handler) CreateSleep(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	var req sleepRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	log, stopped, err := h.Store.CreateSleep(childID, req.StartTime, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	type createSleepResponse struct {
		*model.SleepLog
		StoppedFeeding *store.StoppedFeeding `json:"stopped_feeding,omitempty"`
	}
	h.JSON(w, http.StatusCreated, createSleepResponse{SleepLog: log, StoppedFeeding: stopped})
}

func (h *Handler) GetActiveSleep(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	log, err := h.Store.GetActiveSleep(childID)
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

func (h *Handler) UpdateSleep(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	var req sleepRequest
	if err := h.Decode(r, &req); err != nil {
		h.Error(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	log, err := h.Store.UpdateSleep(id, req.EndTime, req.Notes)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, log)
}

func (h *Handler) DeleteSleep(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	if err := h.Store.DeleteSleep(id); err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
