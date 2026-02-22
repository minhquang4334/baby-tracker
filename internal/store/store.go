package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// hcmcTZ is GMT+7 (Ho Chi Minh City). All timestamps are stored in this zone.
var hcmcTZ = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)

// nowHCMC returns the current time in GMT+7 formatted as RFC3339.
func nowHCMC() string {
	return time.Now().In(hcmcTZ).Format(time.RFC3339)
}

// todayHCMC returns the current date (YYYY-MM-DD) in GMT+7.
func todayHCMC() string {
	return time.Now().In(hcmcTZ).Format("2006-01-02")
}

type Store struct {
	db *sql.DB
}

func Open(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS children (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			date_of_birth TEXT NOT NULL,
			gender TEXT NOT NULL,
			photo_url TEXT DEFAULT '',
			notes TEXT DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sleep_logs (
			id TEXT PRIMARY KEY,
			child_id TEXT NOT NULL REFERENCES children(id),
			start_time TEXT NOT NULL,
			end_time TEXT,
			duration_minutes INTEGER,
			notes TEXT DEFAULT '',
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sleep_child_start ON sleep_logs(child_id, start_time)`,
		`CREATE TABLE IF NOT EXISTS feeding_logs (
			id TEXT PRIMARY KEY,
			child_id TEXT NOT NULL REFERENCES children(id),
			feed_type TEXT NOT NULL,
			start_time TEXT NOT NULL,
			end_time TEXT,
			duration_minutes INTEGER,
			quantity_ml INTEGER,
			notes TEXT DEFAULT '',
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_feeding_child_start ON feeding_logs(child_id, start_time)`,
		`CREATE TABLE IF NOT EXISTS diaper_logs (
			id TEXT PRIMARY KEY,
			child_id TEXT NOT NULL REFERENCES children(id),
			diaper_type TEXT NOT NULL,
			changed_at TEXT NOT NULL,
			notes TEXT DEFAULT '',
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_diaper_child_changed ON diaper_logs(child_id, changed_at)`,
		`CREATE TABLE IF NOT EXISTS growth_logs (
			id TEXT PRIMARY KEY,
			child_id TEXT NOT NULL REFERENCES children(id),
			measured_on TEXT NOT NULL,
			weight_grams INTEGER,
			length_mm INTEGER,
			head_circumference_mm INTEGER,
			notes TEXT DEFAULT '',
			created_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_growth_child_measured ON growth_logs(child_id, measured_on)`,
	}

	for _, stmt := range stmts {
		if _, err := s.db.Exec(stmt); err != nil {
			return fmt.Errorf("exec migration %q: %w", stmt[:40], err)
		}
	}
	return nil
}
