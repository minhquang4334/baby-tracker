package store

import (
	"fmt"
	"time"
)

// DayStats holds aggregated data for a single day.
type DayStats struct {
	Date            string `json:"date"`
	SleepMinutes    int    `json:"sleep_minutes"`
	SleepCount      int    `json:"sleep_count"`
	FeedingCount    int    `json:"feeding_count"`
	BreastFeedCount int    `json:"breast_feed_count"`
	BottleFeedCount int    `json:"bottle_feed_count"`
	BottleMLTotal   int    `json:"bottle_ml_total"`
	DiaperCount     int    `json:"diaper_count"`
	WetCount        int    `json:"wet_count"`
	DirtyCount      int    `json:"dirty_count"`
}

// GetAnalytics returns per-day stats for the child in the [from, to] date range.
func (s *Store) GetAnalytics(childID, from, to string) ([]DayStats, error) {
	stats := map[string]*DayStats{}

	// Sleep aggregation
	sleepRows, err := s.db.Query(`
		SELECT substr(start_time,1,10) as day,
		       COALESCE(SUM(duration_minutes),0),
		       COUNT(*)
		FROM sleep_logs
		WHERE child_id=?
		  AND substr(start_time,1,10) BETWEEN ? AND ?
		  AND end_time IS NOT NULL
		GROUP BY day
		ORDER BY day ASC`, childID, from, to)
	if err != nil {
		return nil, fmt.Errorf("analytics sleep: %w", err)
	}
	defer sleepRows.Close()
	for sleepRows.Next() {
		var day string
		var mins, count int
		if err := sleepRows.Scan(&day, &mins, &count); err != nil {
			return nil, err
		}
		stats[day] = &DayStats{Date: day, SleepMinutes: mins, SleepCount: count}
	}
	if err := sleepRows.Err(); err != nil {
		return nil, err
	}

	// Feeding aggregation
	feedRows, err := s.db.Query(`
		SELECT substr(start_time,1,10) as day,
		       feed_type,
		       COUNT(*),
		       COALESCE(SUM(quantity_ml),0)
		FROM feeding_logs
		WHERE child_id=?
		  AND substr(start_time,1,10) BETWEEN ? AND ?
		GROUP BY day, feed_type
		ORDER BY day ASC`, childID, from, to)
	if err != nil {
		return nil, fmt.Errorf("analytics feeding: %w", err)
	}
	defer feedRows.Close()
	for feedRows.Next() {
		var day, feedType string
		var count, ml int
		if err := feedRows.Scan(&day, &feedType, &count, &ml); err != nil {
			return nil, err
		}
		if stats[day] == nil {
			stats[day] = &DayStats{Date: day}
		}
		d := stats[day]
		d.FeedingCount += count
		if feedType == "bottle" {
			d.BottleFeedCount += count
			d.BottleMLTotal += ml
		} else {
			d.BreastFeedCount += count
		}
	}
	if err := feedRows.Err(); err != nil {
		return nil, err
	}

	// Diaper aggregation
	diaperRows, err := s.db.Query(`
		SELECT substr(changed_at,1,10) as day,
		       diaper_type,
		       COUNT(*)
		FROM diaper_logs
		WHERE child_id=?
		  AND substr(changed_at,1,10) BETWEEN ? AND ?
		GROUP BY day, diaper_type
		ORDER BY day ASC`, childID, from, to)
	if err != nil {
		return nil, fmt.Errorf("analytics diaper: %w", err)
	}
	defer diaperRows.Close()
	for diaperRows.Next() {
		var day, dType string
		var count int
		if err := diaperRows.Scan(&day, &dType, &count); err != nil {
			return nil, err
		}
		if stats[day] == nil {
			stats[day] = &DayStats{Date: day}
		}
		d := stats[day]
		d.DiaperCount += count
		if dType == "wet" {
			d.WetCount += count
		} else if dType == "dirty" {
			d.DirtyCount += count
		}
	}
	if err := diaperRows.Err(); err != nil {
		return nil, err
	}

	return buildDayRange(from, to, stats), nil
}

// buildDayRange produces a slice covering every calendar day in [from, to].
func buildDayRange(from, to string, stats map[string]*DayStats) []DayStats {
	cur, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil
	}
	end, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil
	}

	var out []DayStats
	for !cur.After(end) {
		day := cur.Format("2006-01-02")
		if s, ok := stats[day]; ok {
			out = append(out, *s)
		} else {
			out = append(out, DayStats{Date: day})
		}
		cur = cur.AddDate(0, 0, 1)
	}
	return out
}
