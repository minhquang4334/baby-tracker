package store_test

import (
	"path/filepath"
	"testing"

	"baby-care/internal/store"
)

// newTestStore opens a fresh in-memory-equivalent store backed by a temp file.
func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// mustCreateChild inserts a child and returns its ID.
func mustCreateChild(t *testing.T, st *store.Store) string {
	t.Helper()
	child, err := st.CreateChild("Test Baby", "2024-01-01", "female", "", "")
	if err != nil {
		t.Fatalf("create child: %v", err)
	}
	return child.ID
}

func intPtr(v int) *int { return &v }

func TestOpenAndMigrate(t *testing.T) {
	st := newTestStore(t)
	if st == nil {
		t.Fatal("expected non-nil store")
	}
}
