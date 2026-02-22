package store_test

import (
	"testing"
)

func TestGetDaySummary_Empty(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	summary, err := st.GetDaySummary(childID, "2024-01-15")
	if err != nil {
		t.Fatalf("GetDaySummary: %v", err)
	}
	if summary.SleepCount != 0 {
		t.Errorf("SleepCount = %d, want 0", summary.SleepCount)
	}
	if summary.FeedingCount != 0 {
		t.Errorf("FeedingCount = %d, want 0", summary.FeedingCount)
	}
	if summary.DiaperCount != 0 {
		t.Errorf("DiaperCount = %d, want 0", summary.DiaperCount)
	}
	if summary.LastWeight != nil {
		t.Errorf("LastWeight = %v, want nil", summary.LastWeight)
	}
}

func TestGetDaySummary_Counts(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)
	date := "2024-01-15"

	// 2 completed sleeps on this date
	sl1, _, _ := st.CreateSleep(childID, "2024-01-15T08:00:00+07:00", "")
	st.UpdateSleep(sl1.ID, "", "2024-01-15T09:30:00+07:00", "") // 90 min
	sl2, _, _ := st.CreateSleep(childID, "2024-01-15T13:00:00+07:00", "")
	st.UpdateSleep(sl2.ID, "", "2024-01-15T14:00:00+07:00", "") // 60 min

	// 3 feedings
	st.CreateFeeding(childID, "bottle", "2024-01-15T07:00:00+07:00", "", intPtr(90))
	st.CreateFeeding(childID, "bottle", "2024-01-15T10:00:00+07:00", "", intPtr(120))
	st.CreateFeeding(childID, "bottle", "2024-01-15T16:00:00+07:00", "", intPtr(150))

	// 2 diapers
	st.CreateDiaper(childID, "wet", "2024-01-15T06:00:00+07:00", "")
	st.CreateDiaper(childID, "dirty", "2024-01-15T11:00:00+07:00", "")

	summary, err := st.GetDaySummary(childID, date)
	if err != nil {
		t.Fatalf("GetDaySummary: %v", err)
	}
	if summary.SleepCount != 2 {
		t.Errorf("SleepCount = %d, want 2", summary.SleepCount)
	}
	if summary.TotalSleepMin != 150 {
		t.Errorf("TotalSleepMin = %d, want 150", summary.TotalSleepMin)
	}
	if summary.FeedingCount != 3 {
		t.Errorf("FeedingCount = %d, want 3", summary.FeedingCount)
	}
	if summary.DiaperCount != 2 {
		t.Errorf("DiaperCount = %d, want 2", summary.DiaperCount)
	}
}

func TestGetDaySummary_ExcludesOtherDates(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	sl, _, _ := st.CreateSleep(childID, "2024-01-14T22:00:00+07:00", "")
	st.UpdateSleep(sl.ID, "", "2024-01-14T23:00:00+07:00", "")
	st.CreateFeeding(childID, "bottle", "2024-01-14T20:00:00+07:00", "", intPtr(60))
	st.CreateDiaper(childID, "wet", "2024-01-14T19:00:00+07:00", "")

	// Query a different date
	summary, err := st.GetDaySummary(childID, "2024-01-15")
	if err != nil {
		t.Fatalf("GetDaySummary: %v", err)
	}
	if summary.SleepCount != 0 || summary.FeedingCount != 0 || summary.DiaperCount != 0 {
		t.Errorf("expected all zeros for different date, got sleep=%d feeding=%d diaper=%d",
			summary.SleepCount, summary.FeedingCount, summary.DiaperCount)
	}
}

func TestGetDaySummary_LastWeight(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	st.CreateGrowth(childID, "2024-01-10", intPtr(5000), nil, nil, "")
	st.CreateGrowth(childID, "2024-01-15", intPtr(5200), nil, nil, "")

	summary, err := st.GetDaySummary(childID, "2024-01-15")
	if err != nil {
		t.Fatalf("GetDaySummary: %v", err)
	}
	if summary.LastWeight == nil {
		t.Fatal("expected LastWeight to be set")
	}
	if *summary.LastWeight != 5200 {
		t.Errorf("LastWeight = %d, want 5200", *summary.LastWeight)
	}
}

func TestGetDaySummary_ActiveTimers(t *testing.T) {
	st := newTestStore(t)
	childID := mustCreateChild(t, st)

	sl, _, _ := st.CreateSleep(childID, "2024-01-15T22:00:00+07:00", "")

	summary, err := st.GetDaySummary(childID, "2024-01-15")
	if err != nil {
		t.Fatalf("GetDaySummary: %v", err)
	}
	if summary.ActiveSleep == nil {
		t.Fatal("expected ActiveSleep to be set")
	}
	if summary.ActiveSleep.ID != sl.ID {
		t.Errorf("ActiveSleep.ID = %q, want %q", summary.ActiveSleep.ID, sl.ID)
	}
	// In-progress sleep should not be counted in SleepCount
	if summary.SleepCount != 0 {
		t.Errorf("SleepCount = %d, want 0 (in-progress not counted)", summary.SleepCount)
	}
}
