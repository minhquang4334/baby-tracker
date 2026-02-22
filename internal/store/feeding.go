package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"baby-care/internal/model"
	"github.com/google/uuid"
)

// StoppedSleep is returned by CreateFeeding when an active sleep was auto-stopped.
type StoppedSleep struct {
	ID              string `json:"id"`
	DurationMinutes int    `json:"duration_minutes"`
}

func (s *Store) CreateFeeding(childID, feedType, startTime, notes string, quantityML *int) (*model.FeedingLog, *StoppedSleep, error) {
	now := nowHCMC()
	if startTime == "" {
		startTime = now
	}

	// Auto-stop any active sleep when starting a timed breast feed.
	var stopped *StoppedSleep
	if feedType == "breast_left" || feedType == "breast_right" {
		if activeSleep, err := s.GetActiveSleep(childID); err == nil {
			updated, err := s.UpdateSleep(activeSleep.ID, nowHCMC(), activeSleep.Notes)
			if err == nil && updated.DurationMinutes != nil {
				stopped = &StoppedSleep{ID: updated.ID, DurationMinutes: *updated.DurationMinutes}
			}
		}
	}
	log := &model.FeedingLog{
		ID:         uuid.NewString(),
		ChildID:    childID,
		FeedType:   feedType,
		StartTime:  startTime,
		QuantityML: quantityML,
		Notes:      notes,
		CreatedAt:  now,
	}
	_, err := s.db.Exec(
		`INSERT INTO feeding_logs (id, child_id, feed_type, start_time, quantity_ml, notes, created_at) VALUES (?,?,?,?,?,?,?)`,
		log.ID, log.ChildID, log.FeedType, log.StartTime, log.QuantityML, log.Notes, log.CreatedAt,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("insert feeding: %w", err)
	}
	return log, stopped, nil
}

func (s *Store) GetFeedingLogs(childID, date string) ([]*model.FeedingLog, error) {
	query := `SELECT id, child_id, feed_type, start_time, end_time, duration_minutes, quantity_ml, notes, created_at FROM feeding_logs WHERE child_id=?`
	args := []any{childID}
	if date != "" {
		query += ` AND substr(start_time,1,10)=?`
		args = append(args, date)
	}
	query += ` ORDER BY start_time DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query feeding: %w", err)
	}
	defer rows.Close()
	return scanFeedingRows(rows)
}

func (s *Store) GetActiveFeeding(childID string) (*model.FeedingLog, error) {
	row := s.db.QueryRow(
		`SELECT id, child_id, feed_type, start_time, end_time, duration_minutes, quantity_ml, notes, created_at FROM feeding_logs WHERE child_id=? AND end_time IS NULL AND feed_type != 'bottle' ORDER BY start_time DESC LIMIT 1`,
		childID,
	)
	return scanFeedingRow(row)
}

func (s *Store) UpdateFeeding(id, endTime, notes string, quantityML *int) (*model.FeedingLog, error) {
	var durationMinutes *int
	if endTime != "" {
		var start string
		if err := s.db.QueryRow(`SELECT start_time FROM feeding_logs WHERE id=?`, id).Scan(&start); err == nil {
			st, e1 := time.Parse(time.RFC3339, start)
			et, e2 := time.Parse(time.RFC3339, endTime)
			if e1 == nil && e2 == nil {
				d := int(et.Sub(st).Minutes())
				durationMinutes = &d
			}
		}
		_, err := s.db.Exec(
			`UPDATE feeding_logs SET end_time=?, duration_minutes=?, quantity_ml=?, notes=? WHERE id=?`,
			endTime, durationMinutes, quantityML, notes, id,
		)
		if err != nil {
			return nil, fmt.Errorf("update feeding: %w", err)
		}
	} else {
		_, err := s.db.Exec(`UPDATE feeding_logs SET quantity_ml=?, notes=? WHERE id=?`, quantityML, notes, id)
		if err != nil {
			return nil, fmt.Errorf("update feeding: %w", err)
		}
	}
	return getFeedingByID(s, id)
}

func (s *Store) DeleteFeeding(id string) error {
	_, err := s.db.Exec(`DELETE FROM feeding_logs WHERE id=?`, id)
	return err
}

func getFeedingByID(s *Store, id string) (*model.FeedingLog, error) {
	row := s.db.QueryRow(
		`SELECT id, child_id, feed_type, start_time, end_time, duration_minutes, quantity_ml, notes, created_at FROM feeding_logs WHERE id=?`, id,
	)
	return scanFeedingRow(row)
}

func scanFeedingRow(row *sql.Row) (*model.FeedingLog, error) {
	var l model.FeedingLog
	err := row.Scan(&l.ID, &l.ChildID, &l.FeedType, &l.StartTime, &l.EndTime, &l.DurationMinutes, &l.QuantityML, &l.Notes, &l.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func scanFeedingRows(rows *sql.Rows) ([]*model.FeedingLog, error) {
	var logs []*model.FeedingLog
	for rows.Next() {
		var l model.FeedingLog
		if err := rows.Scan(&l.ID, &l.ChildID, &l.FeedType, &l.StartTime, &l.EndTime, &l.DurationMinutes, &l.QuantityML, &l.Notes, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, rows.Err()
}
