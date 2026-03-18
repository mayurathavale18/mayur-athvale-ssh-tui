package analytics

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const timeFormat = "2006-01-02 15:04:05"

type Visit struct {
	ID             int64          `db:"id"`
	IP             string         `db:"ip"`
	ClientID       string         `db:"client_id"`
	ConnectedAt    string         `db:"connected_at"`
	DisconnectedAt sql.NullString `db:"disconnected_at"`
	Duration       sql.NullInt64  `db:"duration_secs"`
}

type Stats struct {
	TotalVisits   int     `db:"total_visits"`
	UniqueIPs     int     `db:"unique_ips"`
	AvgDuration   float64 `db:"avg_duration"`
	TotalDuration int64   `db:"total_duration"`
}

type RecentVisit struct {
	IP          string `db:"ip"`
	ClientID    string `db:"client_id"`
	ConnectedAt string `db:"connected_at"`
	Duration    string `db:"duration"`
}

type Store struct {
	db *sqlx.DB
}

func NewStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating db directory: %w", err)
	}

	db, err := sqlx.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite handles one writer at a time

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return s, nil
}

func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS visits (
		id              INTEGER PRIMARY KEY AUTOINCREMENT,
		ip              TEXT NOT NULL,
		client_id       TEXT NOT NULL DEFAULT '',
		connected_at    TEXT NOT NULL,
		disconnected_at TEXT,
		duration_secs   INTEGER
	);

	CREATE INDEX IF NOT EXISTS idx_visits_ip ON visits(ip);
	CREATE INDEX IF NOT EXISTS idx_visits_connected_at ON visits(connected_at);
	`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) RecordConnect(ip, clientID string) (int64, error) {
	result, err := s.db.Exec(
		`INSERT INTO visits (ip, client_id, connected_at) VALUES (?, ?, ?)`,
		ip, clientID, time.Now().UTC().Format(timeFormat),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Store) RecordDisconnect(visitID int64) error {
	now := time.Now().UTC().Format(timeFormat)
	_, err := s.db.Exec(
		`UPDATE visits
		 SET disconnected_at = ?,
		     duration_secs = CAST((julianday(?) - julianday(connected_at)) * 86400 AS INTEGER)
		 WHERE id = ?`,
		now, now, visitID,
	)
	return err
}

func (s *Store) GetStats() (Stats, error) {
	var stats Stats
	err := s.db.Get(&stats, `
		SELECT
			COUNT(*) as total_visits,
			COUNT(DISTINCT ip) as unique_ips,
			COALESCE(AVG(duration_secs), 0) as avg_duration,
			COALESCE(SUM(duration_secs), 0) as total_duration
		FROM visits
	`)
	return stats, err
}

func (s *Store) GetRecentVisits(limit int) ([]RecentVisit, error) {
	var visits []RecentVisit
	err := s.db.Select(&visits, `
		SELECT
			ip,
			client_id,
			substr(connected_at, 1, 16) as connected_at,
			CASE
				WHEN duration_secs IS NULL THEN 'active'
				WHEN duration_secs < 60 THEN duration_secs || 's'
				ELSE (duration_secs / 60) || 'm ' || (duration_secs % 60) || 's'
			END as duration
		FROM visits
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	return visits, err
}

func (s *Store) Close() error {
	return s.db.Close()
}
