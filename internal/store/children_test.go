package store_test

import (
	"errors"
	"testing"

	"baby-care/internal/store"
)

func TestGetChild_Empty(t *testing.T) {
	st := newTestStore(t)
	_, err := st.GetChild()
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreateChild(t *testing.T) {
	st := newTestStore(t)
	child, err := st.CreateChild("Lan", "2024-03-15", "female", "", "test notes")
	if err != nil {
		t.Fatalf("CreateChild: %v", err)
	}
	if child.ID == "" {
		t.Error("expected non-empty ID")
	}
	if child.Name != "Lan" {
		t.Errorf("Name = %q, want %q", child.Name, "Lan")
	}
	if child.DateOfBirth != "2024-03-15" {
		t.Errorf("DOB = %q, want %q", child.DateOfBirth, "2024-03-15")
	}
	if child.Gender != "female" {
		t.Errorf("Gender = %q, want %q", child.Gender, "female")
	}
	if child.CreatedAt == "" || child.UpdatedAt == "" {
		t.Error("expected timestamps to be set")
	}
}

func TestGetChild_AfterCreate(t *testing.T) {
	st := newTestStore(t)
	created, _ := st.CreateChild("Minh", "2023-06-01", "male", "", "")

	got, err := st.GetChild()
	if err != nil {
		t.Fatalf("GetChild: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %q, want %q", got.ID, created.ID)
	}
	if got.Name != "Minh" {
		t.Errorf("Name = %q, want %q", got.Name, "Minh")
	}
}

func TestUpdateChild(t *testing.T) {
	st := newTestStore(t)
	child, _ := st.CreateChild("Old Name", "2024-01-01", "male", "", "")

	updated, err := st.UpdateChild(child.ID, "New Name", "2024-02-02", "female", "http://photo.url", "notes")
	if err != nil {
		t.Fatalf("UpdateChild: %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("Name = %q, want %q", updated.Name, "New Name")
	}
	if updated.DateOfBirth != "2024-02-02" {
		t.Errorf("DOB = %q, want %q", updated.DateOfBirth, "2024-02-02")
	}
	if updated.Gender != "female" {
		t.Errorf("Gender = %q, want %q", updated.Gender, "female")
	}
	if updated.PhotoURL != "http://photo.url" {
		t.Errorf("PhotoURL = %q, want %q", updated.PhotoURL, "http://photo.url")
	}
}
