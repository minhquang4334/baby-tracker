package store

import "baby-care/internal/model"

type DaySummary struct {
	Date             string             `json:"date"`
	TotalSleepMin    int                `json:"total_sleep_minutes"`
	SleepCount       int                `json:"sleep_count"`
	FeedingCount     int                `json:"feeding_count"`
	DiaperCount      int                `json:"diaper_count"`
	LastWeight       *int               `json:"last_weight_grams,omitempty"`
	LastSleepEndTime *string            `json:"last_sleep_end_time,omitempty"`
	ActiveSleep      *model.SleepLog   `json:"active_sleep,omitempty"`
	ActiveFeeding    *model.FeedingLog `json:"active_feeding,omitempty"`
}

func (s *Store) GetDaySummary(childID, date string) (*DaySummary, error) {
	summary := &DaySummary{Date: date}

	// Sleep stats
	row := s.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(duration_minutes),0) FROM sleep_logs WHERE child_id=? AND substr(start_time,1,10)=? AND end_time IS NOT NULL`,
		childID, date,
	)
	row.Scan(&summary.SleepCount, &summary.TotalSleepMin)

	// Feeding count
	row = s.db.QueryRow(
		`SELECT COUNT(*) FROM feeding_logs WHERE child_id=? AND substr(start_time,1,10)=?`,
		childID, date,
	)
	row.Scan(&summary.FeedingCount)

	// Diaper count
	row = s.db.QueryRow(
		`SELECT COUNT(*) FROM diaper_logs WHERE child_id=? AND substr(changed_at,1,10)=?`,
		childID, date,
	)
	row.Scan(&summary.DiaperCount)

	// Last weight
	var w int
	err := s.db.QueryRow(
		`SELECT weight_grams FROM growth_logs WHERE child_id=? AND weight_grams IS NOT NULL ORDER BY measured_on DESC LIMIT 1`,
		childID,
	).Scan(&w)
	if err == nil {
		summary.LastWeight = &w
	}

	// Last sleep end time (for awake-time counter)
	var lastSleepEnd string
	err = s.db.QueryRow(
		`SELECT end_time FROM sleep_logs WHERE child_id=? AND end_time IS NOT NULL ORDER BY end_time DESC LIMIT 1`,
		childID,
	).Scan(&lastSleepEnd)
	if err == nil {
		summary.LastSleepEndTime = &lastSleepEnd
	}

	// Active timers
	activeSleep, err := s.GetActiveSleep(childID)
	if err == nil {
		summary.ActiveSleep = activeSleep
	}
	activeFeeding, err := s.GetActiveFeeding(childID)
	if err == nil {
		summary.ActiveFeeding = activeFeeding
	}

	return summary, nil
}
