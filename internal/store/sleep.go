package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"baby-care/internal/model"
	"github.com/google/uuid"
)

// StoppedFeeding is returned by CreateSleep when an active feeding was auto-stopped.
type StoppedFeeding struct {
	ID              string `json:"id"`
	FeedType        string `json:"feed_type"`
	DurationMinutes int    `json:"duration_minutes"`
}

func (s *Store) CreateSleep(childID, startTime, notes string) (*model.SleepLog, *StoppedFeeding, error) {
	now := nowHCMC()
	if startTime == "" {
		startTime = now
	}

	// Auto-stop any active breast feeding session when sleep starts.
	var stopped *StoppedFeeding
	if activeFeeding, err := s.GetActiveFeeding(childID); err == nil {
		updated, err := s.UpdateFeeding(activeFeeding.ID, nowHCMC(), activeFeeding.Notes, activeFeeding.QuantityML)
		if err == nil && updated.DurationMinutes != nil {
			stopped = &StoppedFeeding{ID: updated.ID, FeedType: updated.FeedType, DurationMinutes: *updated.DurationMinutes}
		}
	}

	log := &model.SleepLog{
		ID:        uuid.NewString(),
		ChildID:   childID,
		StartTime: startTime,
		Notes:     notes,
		CreatedAt: now,
	}
	_, err := s.db.Exec(
		`INSERT INTO sleep_logs (id, child_id, start_time, notes, created_at) VALUES (?,?,?,?,?)`,
		log.ID, log.ChildID, log.StartTime, log.Notes, log.CreatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("insert sleep: %w", err)
	}
	return log, stopped, nil
}

func (s *Store) GetSleepLogs(childID, date string) ([]*model.SleepLog, error) {
	query := `SELECT id, child_id, start_time, end_time, duration_minutes, notes, created_at FROM sleep_logs WHERE child_id=?`
	args := []any{childID}
	if date != "" {
		query += ` AND substr(start_time,1,10)=?`
		args = append(args, date)
	}
	query += ` ORDER BY start_time DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query sleep: %w", err)
	}
	defer rows.Close()
	return scanSleepRows(rows)
}

func (s *Store) GetActiveSleep(childID string) (*model.SleepLog, error) {
	row := s.db.QueryRow(
		`SELECT id, child_id, start_time, end_time, duration_minutes, notes, created_at FROM sleep_logs WHERE child_id=? AND end_time IS NULL ORDER BY start_time DESC LIMIT 1`,
		childID,
	)
	return scanSleepRow(row)
}

func (s *Store) UpdateSleep(id, startTime, endTime, notes string) (*model.SleepLog, error) {
	existing, err := getSleepByID(s, id)
	if err != nil {
		return nil, err
	}

	effectiveStart := existing.StartTime
	if startTime != "" {
		effectiveStart = startTime
	}

	var durationMinutes *int
	if existing.EndTime != nil && *existing.EndTime != "" {
		effectiveEnd := *existing.EndTime
		if endTime != "" {
			effectiveEnd = endTime
		}
		st, e1 := time.Parse(time.RFC3339, effectiveStart)
		et, e2 := time.Parse(time.RFC3339, effectiveEnd)
		if e1 == nil && e2 == nil {
			d := int(et.Sub(st).Minutes())
			durationMinutes = &d
		}
	} else if endTime != "" {
		st, e1 := time.Parse(time.RFC3339, effectiveStart)
		et, e2 := time.Parse(time.RFC3339, endTime)
		if e1 == nil && e2 == nil {
			d := int(et.Sub(st).Minutes())
			durationMinutes = &d
		}
	}

	if endTime != "" {
		_, err = s.db.Exec(
			`UPDATE sleep_logs SET start_time=?, end_time=?, duration_minutes=?, notes=? WHERE id=?`,
			effectiveStart, endTime, durationMinutes, notes, id,
		)
	} else {
		_, err = s.db.Exec(
			`UPDATE sleep_logs SET start_time=?, duration_minutes=?, notes=? WHERE id=?`,
			effectiveStart, durationMinutes, notes, id,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("update sleep: %w", err)
	}
	return getSleepByID(s, id)
}

func (s *Store) DeleteSleep(id string) error {
	_, err := s.db.Exec(`DELETE FROM sleep_logs WHERE id=?`, id)
	return err
}

func getSleepByID(s *Store, id string) (*model.SleepLog, error) {
	row := s.db.QueryRow(
		`SELECT id, child_id, start_time, end_time, duration_minutes, notes, created_at FROM sleep_logs WHERE id=?`, id,
	)
	return scanSleepRow(row)
}

func scanSleepRow(row *sql.Row) (*model.SleepLog, error) {
	var l model.SleepLog
	err := row.Scan(&l.ID, &l.ChildID, &l.StartTime, &l.EndTime, &l.DurationMinutes, &l.Notes, &l.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func scanSleepRows(rows *sql.Rows) ([]*model.SleepLog, error) {
	var logs []*model.SleepLog
	for rows.Next() {
		var l model.SleepLog
		if err := rows.Scan(&l.ID, &l.ChildID, &l.StartTime, &l.EndTime, &l.DurationMinutes, &l.Notes, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, rows.Err()
}
