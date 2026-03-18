package analytics

import (
	"sync"
	"testing"
)

func newTestTracker(t *testing.T) *Tracker {
	t.Helper()
	store := newTestStore(t)
	return NewTracker(store)
}

func TestTracker_ActiveCount(t *testing.T) {
	tracker := newTestTracker(t)

	if tracker.ActiveCount() != 0 {
		t.Fatal("expected 0 active initially")
	}

	v1 := tracker.OnConnect("1.1.1.1", "user1")
	v2 := tracker.OnConnect("2.2.2.2", "user2")

	if tracker.ActiveCount() != 2 {
		t.Fatalf("expected 2 active, got %d", tracker.ActiveCount())
	}

	tracker.OnDisconnect(v1)
	if tracker.ActiveCount() != 1 {
		t.Fatalf("expected 1 active after disconnect, got %d", tracker.ActiveCount())
	}

	tracker.OnDisconnect(v2)
	if tracker.ActiveCount() != 0 {
		t.Fatalf("expected 0 active after all disconnect, got %d", tracker.ActiveCount())
	}
}

func TestTracker_ConcurrentConnections(t *testing.T) {
	tracker := newTestTracker(t)

	const n = 50
	var wg sync.WaitGroup
	visitIDs := make([]int64, n)

	// Concurrent connects
	wg.Add(n)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			visitIDs[i] = tracker.OnConnect("10.0.0.1", "user")
		}(i)
	}
	wg.Wait()

	if tracker.ActiveCount() != n {
		t.Fatalf("expected %d active, got %d", n, tracker.ActiveCount())
	}

	// Concurrent disconnects
	wg.Add(n)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			tracker.OnDisconnect(visitIDs[i])
		}(i)
	}
	wg.Wait()

	if tracker.ActiveCount() != 0 {
		t.Fatalf("expected 0 active after all disconnects, got %d", tracker.ActiveCount())
	}
}

func TestTracker_GetStats(t *testing.T) {
	tracker := newTestTracker(t)

	tracker.OnConnect("1.1.1.1", "a")
	tracker.OnConnect("2.2.2.2", "b")

	stats, err := tracker.GetStats()
	if err != nil {
		t.Fatalf("GetStats: %v", err)
	}
	if stats.TotalVisits != 2 {
		t.Fatalf("expected 2 visits, got %d", stats.TotalVisits)
	}
}

func TestTracker_DisconnectRecordsDuration(t *testing.T) {
	tracker := newTestTracker(t)

	id := tracker.OnConnect("1.1.1.1", "user")
	tracker.OnDisconnect(id)

	visits, err := tracker.GetRecentVisits(1)
	if err != nil {
		t.Fatalf("GetRecentVisits: %v", err)
	}
	if len(visits) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(visits))
	}
	if visits[0].Duration == "active" {
		t.Fatal("expected duration to be recorded after disconnect")
	}
}
