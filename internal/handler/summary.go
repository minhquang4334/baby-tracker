package handler

import (
	"net/http"
	"time"
)

var hcmcTZ = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	childID, ok := h.requireChild(w)
	if !ok {
		return
	}
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().In(hcmcTZ).Format("2006-01-02")
	}
	summary, err := h.Store.GetDaySummary(childID, date)
	if err != nil {
		h.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.JSON(w, http.StatusOK, summary)
}
