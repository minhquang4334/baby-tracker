package store_test

import (
	"errors"
	"testing"

	"baby-care/internal/store"
)

const (
	feed1Start = "2024-01-15T10:00:00+07:00"
	feed1End   = "2024-01-15T10:20:00+07:00" // 20 min later
	feed2Start = "2024-01-16T14:00:00+07:00"
)

func TestCreateFeeding_Bottle(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	ml := 120
	f, stopped, err := st.CreateFeeding(childID, "bottle", feed1Start, "", &ml)
	if err != nil {
		t.Fatalf("CreateFeeding: %v", err)
	}
	if stopped != nil {
		t.Error("bottle feed should not stop sleep")
	}
	if f.FeedType != "bottle" {
		t.Errorf("FeedType = %q, want %q", f.FeedType, "bottle")
	}
	if f.QuantityML == nil || *f.QuantityML != 120 {
		t.Errorf("QuantityML = %v, want 120", f.QuantityML)
	}
}

func TestCreateFeeding_BreastLeft(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	f, _, err := st.CreateFeeding(childID, "breast_left", feed1Start, "notes", nil)
	if err != nil {
		t.Fatalf("CreateFeeding: %v", err)
	}
	if f.FeedType != "breast_left" {
		t.Errorf("FeedType = %q, want breast_left", f.FeedType)
	}
	if f.EndTime != nil {
		t.Error("expected EndTime to be nil for active breast feed")
	}
}

func TestCreateFeeding_BreastRight(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	f, _, err := st.CreateFeeding(childID, "breast_right", feed1Start, "", nil)
	if err != nil {
		t.Fatalf("CreateFeeding: %v", err)
	}
	if f.FeedType != "breast_right" {
		t.Errorf("FeedType = %q, want breast_right", f.FeedType)
	}
}

func TestCreateFeeding_DefaultsStartTime(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	f, _, err := st.CreateFeeding(childID, "bottle", "", "", intPtr(60))
	if err != nil {
		t.Fatalf("CreateFeeding: %v", err)
	}
	if f.StartTime == "" {
		t.Error("expected StartTime to be defaulted to now")
	}
}

func TestGetFeedingLogs_All(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateFeeding(childID, "bottle", feed1Start, "", intPtr(90))
	st.CreateFeeding(childID, "breast_left", feed2Start, "", nil)

	logs, err := st.GetFeedingLogs(childID, "")
	if err != nil {
		t.Fatalf("GetFeedingLogs: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("got %d logs, want 2", len(logs))
	}
}

func TestGetFeedingLogs_FilterByDate(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateFeeding(childID, "bottle", feed1Start, "", intPtr(90)) // 2024-01-15
	st.CreateFeeding(childID, "bottle", feed2Start, "", intPtr(60)) // 2024-01-16

	logs, err := st.GetFeedingLogs(childID, "2024-01-15")
	if err != nil {
		t.Fatalf("GetFeedingLogs: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("got %d logs, want 1", len(logs))
	}
}

func TestGetActiveFeeding_NotFound(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	_, err := st.GetActiveFeeding(childID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetActiveFeeding_BreastOnly(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	// Bottle feeds should NOT appear as active
	st.CreateFeeding(childID, "bottle", feed1Start, "", intPtr(90))
	_, err := st.GetActiveFeeding(childID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Error("bottle feed should not be returned as active feeding")
	}

	// Breast feed should appear as active
	breast, _, _ := st.CreateFeeding(childID, "breast_left", feed2Start, "", nil)
	active, err := st.GetActiveFeeding(childID)
	if err != nil {
		t.Fatalf("GetActiveFeeding: %v", err)
	}
	if active.ID != breast.ID {
		t.Errorf("ID = %q, want %q", active.ID, breast.ID)
	}
}

func TestUpdateFeeding_StopsWithDuration(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	f, _, _ := st.CreateFeeding(childID, "breast_left", feed1Start, "", nil)
	updated, err := st.UpdateFeeding(f.ID, "", "", feed1End, "", nil)
	if err != nil {
		t.Fatalf("UpdateFeeding: %v", err)
	}
	if updated.EndTime == nil || *updated.EndTime != feed1End {
		t.Errorf("EndTime = %v, want %q", updated.EndTime, feed1End)
	}
	if updated.DurationMinutes == nil {
		t.Fatal("expected DurationMinutes to be set")
	}
	if *updated.DurationMinutes != 20 {
		t.Errorf("DurationMinutes = %d, want 20", *updated.DurationMinutes)
	}
}

func TestUpdateFeeding_QuantityOnly(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	f, _, _ := st.CreateFeeding(childID, "bottle", feed1Start, "", intPtr(60))
	updated, err := st.UpdateFeeding(f.ID, "", "", "", "", intPtr(150))
	if err != nil {
		t.Fatalf("UpdateFeeding: %v", err)
	}
	if updated.QuantityML == nil || *updated.QuantityML != 150 {
		t.Errorf("QuantityML = %v, want 150", updated.QuantityML)
	}
	if updated.EndTime != nil {
		t.Error("expected EndTime to remain nil")
	}
}

func TestDeleteFeeding(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	f, _, _ := st.CreateFeeding(childID, "bottle", feed1Start, "", intPtr(60))
	if err := st.DeleteFeeding(f.ID); err != nil {
		t.Fatalf("DeleteFeeding: %v", err)
	}

	logs, _ := st.GetFeedingLogs(childID, "")
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}

func TestCreateFeeding_BreastAutoStopsActiveSleep(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	// Start sleep
	sl, _, _ := st.CreateSleep(childID, "2024-01-15T07:00:00+07:00", "")

	// Start breast feed — should auto-stop the sleep
	_, stopped, err := st.CreateFeeding(childID, "breast_left", feed1Start, "", nil)
	if err != nil {
		t.Fatalf("CreateFeeding: %v", err)
	}
	if stopped == nil {
		t.Fatal("expected stopped sleep info")
	}
	if stopped.ID != sl.ID {
		t.Errorf("StoppedSleep.ID = %q, want %q", stopped.ID, sl.ID)
	}
	if stopped.DurationMinutes <= 0 {
		t.Errorf("expected positive duration, got %d", stopped.DurationMinutes)
	}

	// Active sleep should now be gone
	_, err = st.GetActiveSleep(childID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected no active sleep after auto-stop, got %v", err)
	}
}

func TestCreateFeeding_BottleDoesNotAutoStopSleep(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateSleep(childID, "2024-01-15T07:00:00+07:00", "")

	_, stopped, err := st.CreateFeeding(childID, "bottle", feed1Start, "", intPtr(90))
	if err != nil {
		t.Fatalf("CreateFeeding: %v", err)
	}
	if stopped != nil {
		t.Error("bottle feed should not stop active sleep")
	}

	// Sleep should still be active
	_, err = st.GetActiveSleep(childID)
	if err != nil {
		t.Errorf("expected sleep to still be active, got %v", err)
	}
}
