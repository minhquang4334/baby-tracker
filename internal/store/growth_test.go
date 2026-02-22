package store_test

import (
	"testing"
)

func TestCreateGrowth(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	g, err := st.CreateGrowth(childID, "2024-01-15", intPtr(5200), intPtr(580), intPtr(380), "healthy")
	if err != nil {
		t.Fatalf("CreateGrowth: %v", err)
	}
	if g.ID == "" {
		t.Error("expected non-empty ID")
	}
	if g.MeasuredOn != "2024-01-15" {
		t.Errorf("MeasuredOn = %q, want %q", g.MeasuredOn, "2024-01-15")
	}
	if g.WeightGrams == nil || *g.WeightGrams != 5200 {
		t.Errorf("WeightGrams = %v, want 5200", g.WeightGrams)
	}
	if g.LengthMM == nil || *g.LengthMM != 580 {
		t.Errorf("LengthMM = %v, want 580", g.LengthMM)
	}
	if g.HeadCircumferenceMM == nil || *g.HeadCircumferenceMM != 380 {
		t.Errorf("HeadCircumferenceMM = %v, want 380", g.HeadCircumferenceMM)
	}
}

func TestCreateGrowth_DefaultsDate(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	g, err := st.CreateGrowth(childID, "", intPtr(5000), nil, nil, "")
	if err != nil {
		t.Fatalf("CreateGrowth: %v", err)
	}
	if g.MeasuredOn == "" {
		t.Error("expected MeasuredOn to default to today")
	}
}

func TestCreateGrowth_NilFields(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	g, err := st.CreateGrowth(childID, "2024-01-15", nil, nil, nil, "")
	if err != nil {
		t.Fatalf("CreateGrowth with nil fields: %v", err)
	}
	if g.WeightGrams != nil {
		t.Errorf("expected nil WeightGrams, got %v", g.WeightGrams)
	}
}

func TestGetGrowthLogs_OrderedByDate(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateGrowth(childID, "2024-03-01", intPtr(6000), nil, nil, "")
	st.CreateGrowth(childID, "2024-01-01", intPtr(5000), nil, nil, "")
	st.CreateGrowth(childID, "2024-02-01", intPtr(5500), nil, nil, "")

	logs, err := st.GetGrowthLogs(childID)
	if err != nil {
		t.Fatalf("GetGrowthLogs: %v", err)
	}
	if len(logs) != 3 {
		t.Fatalf("got %d logs, want 3", len(logs))
	}
	// Should be ascending by date
	if logs[0].MeasuredOn != "2024-01-01" {
		t.Errorf("first log date = %q, want 2024-01-01", logs[0].MeasuredOn)
	}
	if logs[2].MeasuredOn != "2024-03-01" {
		t.Errorf("last log date = %q, want 2024-03-01", logs[2].MeasuredOn)
	}
}

func TestUpdateGrowth(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	g, _ := st.CreateGrowth(childID, "2024-01-15", intPtr(5000), nil, nil, "")
	updated, err := st.UpdateGrowth(g.ID, "2024-01-16", intPtr(5100), intPtr(590), intPtr(385), "updated")
	if err != nil {
		t.Fatalf("UpdateGrowth: %v", err)
	}
	if updated.MeasuredOn != "2024-01-16" {
		t.Errorf("MeasuredOn = %q, want 2024-01-16", updated.MeasuredOn)
	}
	if updated.WeightGrams == nil || *updated.WeightGrams != 5100 {
		t.Errorf("WeightGrams = %v, want 5100", updated.WeightGrams)
	}
	if updated.Notes != "updated" {
		t.Errorf("Notes = %q, want updated", updated.Notes)
	}
}

func TestDeleteGrowth(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	g, _ := st.CreateGrowth(childID, "2024-01-15", intPtr(5000), nil, nil, "")
	if err := st.DeleteGrowth(g.ID); err != nil {
		t.Fatalf("DeleteGrowth: %v", err)
	}

	logs, _ := st.GetGrowthLogs(childID)
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}
