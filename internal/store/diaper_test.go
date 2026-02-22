package store_test

import (
	"testing"
)

const (
	diaper1Time = "2024-01-15T06:00:00+07:00"
	diaper2Time = "2024-01-16T08:00:00+07:00"
)

func TestCreateDiaper(t *testing.T) {
	tests := []struct {
		name       string
		diaperType string
	}{
		{"wet", "wet"},
		{"dirty", "dirty"},
		{"mixed", "mixed"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := newTestStore(t)
			childID := mustCreateChild(t, st)

			d, err := st.CreateDiaper(childID, tc.diaperType, diaper1Time, "notes")
			if err != nil {
				t.Fatalf("CreateDiaper: %v", err)
			}
			if d.ID == "" {
				t.Error("expected non-empty ID")
			}
			if d.DiaperType != tc.diaperType {
				t.Errorf("DiaperType = %q, want %q", d.DiaperType, tc.diaperType)
			}
			if d.ChangedAt != diaper1Time {
				t.Errorf("ChangedAt = %q, want %q", d.ChangedAt, diaper1Time)
			}
		})
	}
}

func TestCreateDiaper_DefaultsChangedAt(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	d, err := st.CreateDiaper(childID, "wet", "", "")
	if err != nil {
		t.Fatalf("CreateDiaper: %v", err)
	}
	if d.ChangedAt == "" {
		t.Error("expected ChangedAt to be defaulted to now")
	}
}

func TestGetDiaperLogs_All(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateDiaper(childID, "wet", diaper1Time, "")
	st.CreateDiaper(childID, "dirty", diaper2Time, "")

	logs, err := st.GetDiaperLogs(childID, "")
	if err != nil {
		t.Fatalf("GetDiaperLogs: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("got %d logs, want 2", len(logs))
	}
}

func TestGetDiaperLogs_FilterByDate(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateDiaper(childID, "wet", diaper1Time, "")   // 2024-01-15
	st.CreateDiaper(childID, "dirty", diaper2Time, "") // 2024-01-16

	logs, err := st.GetDiaperLogs(childID, "2024-01-15")
	if err != nil {
		t.Fatalf("GetDiaperLogs: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("got %d logs, want 1", len(logs))
	}
	if logs[0].DiaperType != "wet" {
		t.Errorf("DiaperType = %q, want wet", logs[0].DiaperType)
	}
}

func TestUpdateDiaper(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	d, _ := st.CreateDiaper(childID, "wet", diaper1Time, "")
	updated, err := st.UpdateDiaper(d.ID, "mixed", "updated note")
	if err != nil {
		t.Fatalf("UpdateDiaper: %v", err)
	}
	if updated.DiaperType != "mixed" {
		t.Errorf("DiaperType = %q, want mixed", updated.DiaperType)
	}
	if updated.Notes != "updated note" {
		t.Errorf("Notes = %q, want %q", updated.Notes, "updated note")
	}
}

func TestDeleteDiaper(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	d, _ := st.CreateDiaper(childID, "wet", diaper1Time, "")
	if err := st.DeleteDiaper(d.ID); err != nil {
		t.Fatalf("DeleteDiaper: %v", err)
	}

	logs, _ := st.GetDiaperLogs(childID, "")
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}
