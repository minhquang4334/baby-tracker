package handler

import (
	"net/http"
	"time"

	"baby-care/internal/store"
)

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}

	to := r.URL.Query().Get("to")
	from := r.URL.Query().Get("from")

	// Default: last 7 days in GMT+7
	now := time.Now().In(hcmcTZ)
	if to == "" {
		to = now.Format("2006-01-02")
	}
	if from == "" {
		from = now.AddDate(0, 0, -6).Format("2006-01-02")
	}

	days, err := h.Store.GetAnalytics(childID, from, to)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if days == nil {
		days = []store.DayStats{}
	}
	h.JSON(w, http.StatusOK, days)
}
