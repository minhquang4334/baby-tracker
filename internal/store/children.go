package store

import (
	"database/sql"
	"errors"
	"fmt"

	"baby-care/internal/model"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

func (s *Store) GetChild() (*model.Child, error) {
	row := s.db.QueryRow(`SELECT id, name, date_of_birth, gender, photo_url, notes, created_at, updated_at FROM children LIMIT 1`)
	var c model.Child
	err := row.Scan(&c.ID, &c.Name, &c.DateOfBirth, &c.Gender, &c.PhotoURL, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan child: %w", err)
	}
	return &c, nil
}

func (s *Store) CreateChild(name, dob, gender, photoURL, notes string) (*model.Child, error) {
	now := nowHCMC()
	c := &model.Child{
		ID:          uuid.NewString(),
		Name:        name,
		DateOfBirth: dob,
		Gender:      gender,
		PhotoURL:    photoURL,
		Notes:       notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := s.db.Exec(
		`INSERT INTO children (id, name, date_of_birth, gender, photo_url, notes, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?)`,
		c.ID, c.Name, c.DateOfBirth, c.Gender, c.PhotoURL, c.Notes, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert child: %w", err)
	}
	return c, nil
}

func (s *Store) UpdateChild(id, name, dob, gender, photoURL, notes string) (*model.Child, error) {
	now := nowHCMC()
	_, err := s.db.Exec(
		`UPDATE children SET name=?, date_of_birth=?, gender=?, photo_url=?, notes=?, updated_at=? WHERE id=?`,
		name, dob, gender, photoURL, notes, now, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update child: %w", err)
	}
	return s.GetChild()
}
