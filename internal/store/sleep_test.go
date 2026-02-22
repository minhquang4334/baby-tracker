package store_test

import (
	"errors"
	"testing"

	"baby-care/internal/store"
)

const (
	sleep1Start = "2024-01-15T08:00:00+07:00"
	sleep1End   = "2024-01-15T09:30:00+07:00" // 90 min later
	sleep2Start = "2024-01-16T21:00:00+07:00"
)

func TestCreateSleep_ExplicitTime(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	sl, stopped, err := st.CreateSleep(childID, sleep1Start, "nap")
	if err != nil {
		t.Fatalf("CreateSleep: %v", err)
	}
	if stopped != nil {
		t.Error("expected no stopped feeding")
	}
	if sl.ID == "" {
		t.Error("expected non-empty ID")
	}
	if sl.StartTime != sleep1Start {
		t.Errorf("StartTime = %q, want %q", sl.StartTime, sleep1Start)
	}
	if sl.EndTime != nil {
		t.Error("expected EndTime to be nil")
	}
	if sl.Notes != "nap" {
		t.Errorf("Notes = %q, want %q", sl.Notes, "nap")
	}
}

func TestCreateSleep_DefaultsStartTime(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	sl, _, err := st.CreateSleep(childID, "", "")
	if err != nil {
		t.Fatalf("CreateSleep: %v", err)
	}
	if sl.StartTime == "" {
		t.Error("expected StartTime to be defaulted to now")
	}
}

func TestGetSleepLogs_All(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateSleep(childID, sleep1Start, "")
	st.CreateSleep(childID, sleep2Start, "")

	logs, err := st.GetSleepLogs(childID, "")
	if err != nil {
		t.Fatalf("GetSleepLogs: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("got %d logs, want 2", len(logs))
	}
}

func TestGetSleepLogs_FilterByDate(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateSleep(childID, sleep1Start, "")  // 2024-01-15
	st.CreateSleep(childID, sleep2Start, "")  // 2024-01-16

	logs, err := st.GetSleepLogs(childID, "2024-01-15")
	if err != nil {
		t.Fatalf("GetSleepLogs: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("got %d logs, want 1", len(logs))
	}
	if logs[0].StartTime != sleep1Start {
		t.Errorf("unexpected log returned")
	}
}

func TestGetActiveSleep_NotFound(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	_, err := st.GetActiveSleep(childID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetActiveSleep_Found(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	created, _, _ := st.CreateSleep(childID, sleep1Start, "")

	active, err := st.GetActiveSleep(childID)
	if err != nil {
		t.Fatalf("GetActiveSleep: %v", err)
	}
	if active.ID != created.ID {
		t.Errorf("ID = %q, want %q", active.ID, created.ID)
	}
}

func TestUpdateSleep_StopsWithDuration(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	sl, _, _ := st.CreateSleep(childID, sleep1Start, "")
	updated, err := st.UpdateSleep(sl.ID, "", sleep1End, "updated notes")
	if err != nil {
		t.Fatalf("UpdateSleep: %v", err)
	}
	if updated.EndTime == nil || *updated.EndTime != sleep1End {
		t.Errorf("EndTime = %v, want %q", updated.EndTime, sleep1End)
	}
	if updated.DurationMinutes == nil {
		t.Fatal("expected DurationMinutes to be set")
	}
	if *updated.DurationMinutes != 90 {
		t.Errorf("DurationMinutes = %d, want 90", *updated.DurationMinutes)
	}
	if updated.Notes != "updated notes" {
		t.Errorf("Notes = %q, want %q", updated.Notes, "updated notes")
	}
}

func TestUpdateSleep_NotesOnly(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	sl, _, _ := st.CreateSleep(childID, sleep1Start, "original")
	updated, err := st.UpdateSleep(sl.ID, "", "", "changed")
	if err != nil {
		t.Fatalf("UpdateSleep: %v", err)
	}
	if updated.EndTime != nil {
		t.Error("expected EndTime to remain nil")
	}
	if updated.Notes != "changed" {
		t.Errorf("Notes = %q, want %q", updated.Notes, "changed")
	}
}

func TestDeleteSleep(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	sl, _, _ := st.CreateSleep(childID, sleep1Start, "")
	if err := st.DeleteSleep(sl.ID); err != nil {
		t.Fatalf("DeleteSleep: %v", err)
	}

	logs, _ := st.GetSleepLogs(childID, "")
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}

func TestCreateSleep_AutoStopsActiveFeeding(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	// Start a breast feed
	feed, _, _ := st.CreateFeeding(childID, "breast_left", "2024-01-15T07:00:00+07:00", "", nil)
	if feed == nil {
		t.Fatal("expected feeding to be created")
	}

	// Start sleep — should auto-stop the active feeding
	_, stopped, err := st.CreateSleep(childID, sleep1Start, "")
	if err != nil {
		t.Fatalf("CreateSleep: %v", err)
	}
	if stopped == nil {
		t.Fatal("expected stopped feeding info")
	}
	if stopped.ID != feed.ID {
		t.Errorf("StoppedFeeding.ID = %q, want %q", stopped.ID, feed.ID)
	}
	if stopped.DurationMinutes <= 0 {
		t.Errorf("expected positive duration, got %d", stopped.DurationMinutes)
	}

	// Active feeding should now be gone
	_, err = st.GetActiveFeeding(childID)
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("expected no active feeding after auto-stop, got %v", err)
	}
}
