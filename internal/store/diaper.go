package store

import (
	"database/sql"
	"errors"
	"fmt"

	"baby-care/internal/model"
	"github.com/google/uuid"
)

func (s *Store) CreateDiaper(childID, diaperType, changedAt, notes string) (*model.DiaperLog, error) {
	now := nowHCMC()
	if changedAt == "" {
		changedAt = now
	}
	log := &model.DiaperLog{
		ID:         uuid.NewString(),
		ChildID:    childID,
		DiaperType: diaperType,
		ChangedAt:  changedAt,
		Notes:      notes,
		CreatedAt:  now,
	}
	_, err := s.db.Exec(
		`INSERT INTO diaper_logs (id, child_id, diaper_type, changed_at, notes, created_at) VALUES (?,?,?,?,?,?)`,
		log.ID, log.ChildID, log.DiaperType, log.ChangedAt, log.Notes, log.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert diaper: %w", err)
	}
	return log, nil
}

func (s *Store) GetDiaperLogs(childID, date string) ([]*model.DiaperLog, error) {
	query := `SELECT id, child_id, diaper_type, changed_at, notes, created_at FROM diaper_logs WHERE child_id=?`
	args := []any{childID}
	if date != "" {
		query += ` AND substr(changed_at,1,10)=?`
		args = append(args, date)
	}
	query += ` ORDER BY changed_at DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query diaper: %w", err)
	}
	defer rows.Close()
	return scanDiaperRows(rows)
}

func (s *Store) UpdateDiaper(id, diaperType, changedAt, notes string) (*model.DiaperLog, error) {
	existing, err := getDiaperByID(s, id)
	if err != nil {
		return nil, err
	}
	if changedAt == "" {
		changedAt = existing.ChangedAt
	}
	_, err = s.db.Exec(
		`UPDATE diaper_logs SET diaper_type=?, changed_at=?, notes=? WHERE id=?`,
		diaperType, changedAt, notes, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update diaper: %w", err)
	}
	return getDiaperByID(s, id)
}

func (s *Store) DeleteDiaper(id string) error {
	_, err := s.db.Exec(`DELETE FROM diaper_logs WHERE id=?`, id)
	return err
}

func getDiaperByID(s *Store, id string) (*model.DiaperLog, error) {
	row := s.db.QueryRow(
		`SELECT id, child_id, diaper_type, changed_at, notes, created_at FROM diaper_logs WHERE id=?`, id,
	)
	var l model.DiaperLog
	err := row.Scan(&l.ID, &l.ChildID, &l.DiaperType, &l.ChangedAt, &l.Notes, &l.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func scanDiaperRows(rows *sql.Rows) ([]*model.DiaperLog, error) {
	var logs []*model.DiaperLog
	for rows.Next() {
		var l model.DiaperLog
		if err := rows.Scan(&l.ID, &l.ChildID, &l.DiaperType, &l.ChangedAt, &l.Notes, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, rows.Err()
}
