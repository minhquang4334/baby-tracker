package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"baby-care/internal/model"
	"baby-care/internal/server"
	"baby-care/internal/store"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func newTestServer(t *testing.T) (*httptest.Server, *store.Store) {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })

	staticFS := fstest.MapFS{
		"index.html": {Data: []byte("<html></html>")},
	}
	srv := httptest.NewServer(server.New(st, staticFS))
	t.Cleanup(srv.Close)
	return srv, st
}

func do(t *testing.T, srv *httptest.Server, method, path string, body any) *http.Response {
	t.Helper()
	var buf *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	} else {
		buf = &bytes.Buffer{}
	}
	req, err := http.NewRequest(method, srv.URL+path, buf)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	return resp
}

func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	resp.Body.Close()
}

func mustCreateChildViaAPI(t *testing.T, srv *httptest.Server) *model.Child {
	t.Helper()
	resp := do(t, srv, "POST", "/api/v1/child", map[string]string{
		"name": "Test Baby", "date_of_birth": "2024-01-01", "gender": "female",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create child: status %d", resp.StatusCode)
	}
	var child model.Child
	decodeJSON(t, resp, &child)
	return &child
}

func intPtr(v int) *int { return &v }

// ── health ────────────────────────────────────────────────────────────────────

func TestHealth(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := do(t, srv, "GET", "/health", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

// ── child ─────────────────────────────────────────────────────────────────────

func TestGetChild_NotFound(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := do(t, srv, "GET", "/api/v1/child", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestCreateChild_MissingFields(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := do(t, srv, "POST", "/api/v1/child", map[string]string{"name": "Only Name"})
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestCreateChild_InvalidJSON(t *testing.T) {
	srv, _ := newTestServer(t)
	req, _ := http.NewRequest("POST", srv.URL+"/api/v1/child", strings.NewReader("not json"))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestCreateAndGetChild(t *testing.T) {
	srv, _ := newTestServer(t)
	child := mustCreateChildViaAPI(t, srv)

	if child.Name != "Test Baby" {
		t.Errorf("Name = %q, want Test Baby", child.Name)
	}
	if child.ID == "" {
		t.Error("expected non-empty ID")
	}

	resp := do(t, srv, "GET", "/api/v1/child", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET child status = %d, want 200", resp.StatusCode)
	}
	var got model.Child
	decodeJSON(t, resp, &got)
	if got.ID != child.ID {
		t.Errorf("ID = %q, want %q", got.ID, child.ID)
	}
}

func TestUpdateChild(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "PUT", "/api/v1/child", map[string]string{
		"name": "Updated Name", "date_of_birth": "2024-06-01", "gender": "male",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	var updated model.Child
	decodeJSON(t, resp, &updated)
	if updated.Name != "Updated Name" {
		t.Errorf("Name = %q, want Updated Name", updated.Name)
	}
}

func TestUpdateChild_NoChild(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := do(t, srv, "PUT", "/api/v1/child", map[string]string{
		"name": "X", "date_of_birth": "2024-01-01", "gender": "male",
	})
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

// ── sleep ─────────────────────────────────────────────────────────────────────

func TestSleep_RequiresChild(t *testing.T) {
	srv, _ := newTestServer(t)
	for _, tc := range []struct{ method, path string }{
		{"GET", "/api/v1/sleep"},
		{"POST", "/api/v1/sleep"},
		{"GET", "/api/v1/sleep/active"},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			resp := do(t, srv, tc.method, tc.path, map[string]string{})
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400 (no child)", resp.StatusCode)
			}
		})
	}
}

func TestCreateAndListSleep(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	// Create a sleep log
	resp := do(t, srv, "POST", "/api/v1/sleep", map[string]string{
		"start_time": "2024-01-15T22:00:00+07:00",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create sleep status = %d, want 201", resp.StatusCode)
	}
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// List sleeps
	resp = do(t, srv, "GET", "/api/v1/sleep", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list sleep status = %d, want 200", resp.StatusCode)
	}
	var logs []map[string]any
	decodeJSON(t, resp, &logs)
	if len(logs) != 1 {
		t.Errorf("got %d logs, want 1", len(logs))
	}
}

func TestGetActiveSleep_None(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "GET", "/api/v1/sleep/active", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	// Body should be null
	var val any
	decodeJSON(t, resp, &val)
	if val != nil {
		t.Errorf("expected null, got %v", val)
	}
}

func TestUpdateSleep_Stop(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/sleep", map[string]string{
		"start_time": "2024-01-15T08:00:00+07:00",
	})
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)

	resp = do(t, srv, "PUT", "/api/v1/sleep/"+id, map[string]string{
		"end_time": "2024-01-15T09:00:00+07:00",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update sleep status = %d, want 200", resp.StatusCode)
	}
	var updated map[string]any
	decodeJSON(t, resp, &updated)
	if updated["end_time"] == nil {
		t.Error("expected end_time to be set")
	}
	if updated["duration_minutes"].(float64) != 60 {
		t.Errorf("duration_minutes = %v, want 60", updated["duration_minutes"])
	}
}

func TestDeleteSleep(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/sleep", map[string]string{
		"start_time": "2024-01-15T08:00:00+07:00",
	})
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)

	resp = do(t, srv, "DELETE", "/api/v1/sleep/"+id, nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("delete status = %d, want 204", resp.StatusCode)
	}

	resp = do(t, srv, "GET", "/api/v1/sleep", nil)
	var logs []any
	decodeJSON(t, resp, &logs)
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}

func TestCreateSleep_AutoStopsFeeding(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	// Start a breast feed
	resp := do(t, srv, "POST", "/api/v1/feeding", map[string]string{
		"feed_type": "breast_left", "start_time": "2024-01-15T07:00:00+07:00",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create feeding status = %d", resp.StatusCode)
	}

	// Start sleep — response should include stopped_feeding
	resp = do(t, srv, "POST", "/api/v1/sleep", map[string]string{
		"start_time": "2024-01-15T08:00:00+07:00",
	})
	var sleepResp map[string]any
	decodeJSON(t, resp, &sleepResp)
	if sleepResp["stopped_feeding"] == nil {
		t.Error("expected stopped_feeding in response")
	}
}

// ── feeding ───────────────────────────────────────────────────────────────────

func TestCreateFeeding_MissingFeedType(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/feeding", map[string]string{
		"start_time": "2024-01-15T10:00:00+07:00",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestCreateFeeding_Bottle(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/feeding", map[string]any{
		"feed_type": "bottle", "start_time": "2024-01-15T10:00:00+07:00", "quantity_ml": 120,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", resp.StatusCode)
	}
	var created map[string]any
	decodeJSON(t, resp, &created)
	if created["feed_type"] != "bottle" {
		t.Errorf("feed_type = %v, want bottle", created["feed_type"])
	}
	if created["quantity_ml"].(float64) != 120 {
		t.Errorf("quantity_ml = %v, want 120", created["quantity_ml"])
	}
}

func TestGetActiveFeeding_None(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "GET", "/api/v1/feeding/active", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var val any
	decodeJSON(t, resp, &val)
	if val != nil {
		t.Errorf("expected null active feeding, got %v", val)
	}
}

func TestUpdateFeeding_Quantity(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/feeding", map[string]any{
		"feed_type": "bottle", "start_time": "2024-01-15T10:00:00+07:00", "quantity_ml": 90,
	})
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)

	resp = do(t, srv, "PUT", "/api/v1/feeding/"+id, map[string]any{"quantity_ml": 150})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update feeding status = %d, want 200", resp.StatusCode)
	}
	var updated map[string]any
	decodeJSON(t, resp, &updated)
	if updated["quantity_ml"].(float64) != 150 {
		t.Errorf("quantity_ml = %v, want 150", updated["quantity_ml"])
	}
}

func TestDeleteFeeding(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/feeding", map[string]any{
		"feed_type": "bottle", "start_time": "2024-01-15T10:00:00+07:00", "quantity_ml": 90,
	})
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)

	resp = do(t, srv, "DELETE", "/api/v1/feeding/"+id, nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("delete status = %d, want 204", resp.StatusCode)
	}

	resp = do(t, srv, "GET", "/api/v1/feeding", nil)
	var logs []any
	decodeJSON(t, resp, &logs)
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}

func TestCreateFeeding_BreastAutoStopsSleep(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	// Start sleep
	do(t, srv, "POST", "/api/v1/sleep", map[string]string{
		"start_time": "2024-01-15T07:00:00+07:00",
	})

	// Start breast feed — should auto-stop sleep
	resp := do(t, srv, "POST", "/api/v1/feeding", map[string]string{
		"feed_type": "breast_right", "start_time": "2024-01-15T08:00:00+07:00",
	})
	var feedResp map[string]any
	decodeJSON(t, resp, &feedResp)
	if feedResp["stopped_sleep"] == nil {
		t.Error("expected stopped_sleep in response")
	}
}

// ── diaper ────────────────────────────────────────────────────────────────────

func TestCreateDiaper_MissingType(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/diaper", map[string]string{
		"changed_at": "2024-01-15T06:00:00+07:00",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
}

func TestCreateAndListDiaper(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	for _, dt := range []string{"wet", "dirty", "mixed"} {
		resp := do(t, srv, "POST", "/api/v1/diaper", map[string]string{
			"diaper_type": dt, "changed_at": "2024-01-15T06:00:00+07:00",
		})
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("create %s diaper: status = %d, want 201", dt, resp.StatusCode)
		}
		resp.Body.Close()
	}

	resp := do(t, srv, "GET", "/api/v1/diaper", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list diaper status = %d, want 200", resp.StatusCode)
	}
	var logs []any
	decodeJSON(t, resp, &logs)
	if len(logs) != 3 {
		t.Errorf("got %d logs, want 3", len(logs))
	}
}

func TestDeleteDiaper(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/diaper", map[string]string{
		"diaper_type": "wet", "changed_at": "2024-01-15T06:00:00+07:00",
	})
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)

	resp = do(t, srv, "DELETE", "/api/v1/diaper/"+id, nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("delete status = %d, want 204", resp.StatusCode)
	}
}

// ── growth ────────────────────────────────────────────────────────────────────

func TestCreateAndListGrowth(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/growth", map[string]any{
		"measured_on": "2024-01-15", "weight_grams": 5200, "length_mm": 580,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create growth status = %d, want 201", resp.StatusCode)
	}
	var created map[string]any
	decodeJSON(t, resp, &created)
	if created["weight_grams"].(float64) != 5200 {
		t.Errorf("weight_grams = %v, want 5200", created["weight_grams"])
	}

	resp = do(t, srv, "GET", "/api/v1/growth", nil)
	var logs []any
	decodeJSON(t, resp, &logs)
	if len(logs) != 1 {
		t.Errorf("got %d logs, want 1", len(logs))
	}
}

func TestUpdateGrowth(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/growth", map[string]any{
		"measured_on": "2024-01-15", "weight_grams": 5000,
	})
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)

	resp = do(t, srv, "PUT", "/api/v1/growth/"+id, map[string]any{
		"measured_on": "2024-01-16", "weight_grams": 5100,
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update growth status = %d, want 200", resp.StatusCode)
	}
	var updated map[string]any
	decodeJSON(t, resp, &updated)
	if updated["weight_grams"].(float64) != 5100 {
		t.Errorf("weight_grams = %v, want 5100", updated["weight_grams"])
	}
}

func TestDeleteGrowth(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "POST", "/api/v1/growth", map[string]any{
		"measured_on": "2024-01-15", "weight_grams": 5000,
	})
	var created map[string]any
	decodeJSON(t, resp, &created)
	id := created["id"].(string)

	resp = do(t, srv, "DELETE", "/api/v1/growth/"+id, nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("delete status = %d, want 204", resp.StatusCode)
	}
}

// ── summary ───────────────────────────────────────────────────────────────────

func TestGetSummary(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	// Add some data for 2024-01-15
	resp := do(t, srv, "POST", "/api/v1/sleep", map[string]string{
		"start_time": "2024-01-15T08:00:00+07:00",
	})
	var sl map[string]any
	decodeJSON(t, resp, &sl)
	do(t, srv, "PUT", "/api/v1/sleep/"+sl["id"].(string), map[string]string{
		"end_time": "2024-01-15T09:00:00+07:00",
	})

	do(t, srv, "POST", "/api/v1/diaper", map[string]string{
		"diaper_type": "wet", "changed_at": "2024-01-15T07:00:00+07:00",
	})

	resp = do(t, srv, "GET", "/api/v1/summary?date=2024-01-15", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("summary status = %d, want 200", resp.StatusCode)
	}
	var summary map[string]any
	decodeJSON(t, resp, &summary)

	if summary["sleep_count"].(float64) != 1 {
		t.Errorf("sleep_count = %v, want 1", summary["sleep_count"])
	}
	if summary["total_sleep_minutes"].(float64) != 60 {
		t.Errorf("total_sleep_minutes = %v, want 60", summary["total_sleep_minutes"])
	}
	if summary["diaper_count"].(float64) != 1 {
		t.Errorf("diaper_count = %v, want 1", summary["diaper_count"])
	}
}

func TestGetSummary_DefaultsToToday(t *testing.T) {
	srv, _ := newTestServer(t)
	mustCreateChildViaAPI(t, srv)

	resp := do(t, srv, "GET", "/api/v1/summary", nil)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

// ── CORS ──────────────────────────────────────────────────────────────────────

func TestCORSHeaders(t *testing.T) {
	srv, _ := newTestServer(t)
	resp := do(t, srv, "GET", "/health", nil)
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header on response")
	}
}
