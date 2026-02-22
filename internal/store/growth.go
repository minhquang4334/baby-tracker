package store

import (
	"database/sql"
	"errors"
	"fmt"

	"baby-care/internal/model"
	"github.com/google/uuid"
)

func (s *Store) CreateGrowth(childID, measuredOn string, weightGrams, lengthMM, headCircMM *int, notes string) (*model.GrowthLog, error) {
	now := nowHCMC()
	if measuredOn == "" {
		measuredOn = todayHCMC()
	}
	log := &model.GrowthLog{
		ID:                  uuid.NewString(),
		ChildID:             childID,
		MeasuredOn:          measuredOn,
		WeightGrams:         weightGrams,
		LengthMM:            lengthMM,
		HeadCircumferenceMM: headCircMM,
		Notes:               notes,
		CreatedAt:           now,
	}
	_, err := s.db.Exec(
		`INSERT INTO growth_logs (id, child_id, measured_on, weight_grams, length_mm, head_circumference_mm, notes, created_at) VALUES (?,?,?,?,?,?,?,?)`,
		log.ID, log.ChildID, log.MeasuredOn, log.WeightGrams, log.LengthMM, log.HeadCircumferenceMM, log.Notes, log.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert growth: %w", err)
	}
	return log, nil
}

func (s *Store) GetGrowthLogs(childID string) ([]*model.GrowthLog, error) {
	rows, err := s.db.Query(
		`SELECT id, child_id, measured_on, weight_grams, length_mm, head_circumference_mm, notes, created_at FROM growth_logs WHERE child_id=? ORDER BY measured_on ASC`,
		childID,
	)
	if err != nil {
		return nil, fmt.Errorf("query growth: %w", err)
	}
	defer rows.Close()
	return scanGrowthRows(rows)
}

func (s *Store) UpdateGrowth(id, measuredOn string, weightGrams, lengthMM, headCircMM *int, notes string) (*model.GrowthLog, error) {
	_, err := s.db.Exec(
		`UPDATE growth_logs SET measured_on=?, weight_grams=?, length_mm=?, head_circumference_mm=?, notes=? WHERE id=?`,
		measuredOn, weightGrams, lengthMM, headCircMM, notes, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update growth: %w", err)
	}
	return getGrowthByID(s, id)
}

func (s *Store) DeleteGrowth(id string) error {
	_, err := s.db.Exec(`DELETE FROM growth_logs WHERE id=?`, id)
	return err
}

func getGrowthByID(s *Store, id string) (*model.GrowthLog, error) {
	row := s.db.QueryRow(
		`SELECT id, child_id, measured_on, weight_grams, length_mm, head_circumference_mm, notes, created_at FROM growth_logs WHERE id=?`, id,
	)
	var l model.GrowthLog
	err := row.Scan(&l.ID, &l.ChildID, &l.MeasuredOn, &l.WeightGrams, &l.LengthMM, &l.HeadCircumferenceMM, &l.Notes, &l.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func scanGrowthRows(rows *sql.Rows) ([]*model.GrowthLog, error) {
	var logs []*model.GrowthLog
	for rows.Next() {
		var l model.GrowthLog
		if err := rows.Scan(&l.ID, &l.ChildID, &l.MeasuredOn, &l.WeightGrams, &l.LengthMM, &l.HeadCircumferenceMM, &l.Notes, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, &l)
	}
	return logs, rows.Err()
}
