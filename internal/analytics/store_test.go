package analytics

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestNewStore_CreatesDBFile(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "subdir", "test.db")

	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected database file to be created")
	}
}

func TestRecordConnect(t *testing.T) {
	store := newTestStore(t)

	id, err := store.RecordConnect("192.168.1.1", "testuser")
	if err != nil {
		t.Fatalf("RecordConnect: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected positive visit ID, got %d", id)
	}

	// Second connection should get a different ID
	id2, err := store.RecordConnect("192.168.1.2", "testuser2")
	if err != nil {
		t.Fatalf("RecordConnect: %v", err)
	}
	if id2 <= id {
		t.Fatalf("expected id2 (%d) > id (%d)", id2, id)
	}
}

func TestRecordDisconnect(t *testing.T) {
	store := newTestStore(t)

	id, _ := store.RecordConnect("10.0.0.1", "user1")

	// Small delay so duration > 0
	time.Sleep(10 * time.Millisecond)

	err := store.RecordDisconnect(id)
	if err != nil {
		t.Fatalf("RecordDisconnect: %v", err)
	}

	// Verify via recent visits
	visits, err := store.GetRecentVisits(1)
	if err != nil {
		t.Fatalf("GetRecentVisits: %v", err)
	}
	if len(visits) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(visits))
	}
	if visits[0].Duration == "active" {
		t.Fatal("expected duration to be set after disconnect, got 'active'")
	}
}

func TestGetStats_Empty(t *testing.T) {
	store := newTestStore(t)

	stats, err := store.GetStats()
	if err != nil {
		t.Fatalf("GetStats: %v", err)
	}
	if stats.TotalVisits != 0 {
		t.Fatalf("expected 0 total visits, got %d", stats.TotalVisits)
	}
	if stats.UniqueIPs != 0 {
		t.Fatalf("expected 0 unique IPs, got %d", stats.UniqueIPs)
	}
}

func TestGetStats_WithData(t *testing.T) {
	store := newTestStore(t)

	// 3 visits from 2 unique IPs
	store.RecordConnect("192.168.1.1", "a")
	store.RecordConnect("192.168.1.1", "b")
	store.RecordConnect("10.0.0.5", "c")

	stats, err := store.GetStats()
	if err != nil {
		t.Fatalf("GetStats: %v", err)
	}
	if stats.TotalVisits != 3 {
		t.Fatalf("expected 3 total visits, got %d", stats.TotalVisits)
	}
	if stats.UniqueIPs != 2 {
		t.Fatalf("expected 2 unique IPs, got %d", stats.UniqueIPs)
	}
}

func TestGetRecentVisits_Order(t *testing.T) {
	store := newTestStore(t)

	store.RecordConnect("1.1.1.1", "first")
	time.Sleep(10 * time.Millisecond)
	store.RecordConnect("2.2.2.2", "second")

	visits, err := store.GetRecentVisits(10)
	if err != nil {
		t.Fatalf("GetRecentVisits: %v", err)
	}
	if len(visits) != 2 {
		t.Fatalf("expected 2 visits, got %d", len(visits))
	}

	// Most recent first
	if visits[0].IP != "2.2.2.2" {
		t.Fatalf("expected most recent visit first, got IP %s", visits[0].IP)
	}
}

func TestGetRecentVisits_Limit(t *testing.T) {
	store := newTestStore(t)

	for i := range 5 {
		store.RecordConnect("10.0.0."+string(rune('1'+i)), "user")
	}

	visits, err := store.GetRecentVisits(3)
	if err != nil {
		t.Fatalf("GetRecentVisits: %v", err)
	}
	if len(visits) != 3 {
		t.Fatalf("expected 3 visits, got %d", len(visits))
	}
}
