package analytics

import (
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/log"
)

type Tracker struct {
	store       *Store
	activeCount atomic.Int32
	mu          sync.Mutex
}

func NewTracker(store *Store) *Tracker {
	return &Tracker{store: store}
}

func (t *Tracker) OnConnect(ip, clientID string) int64 {
	t.activeCount.Add(1)

	visitID, err := t.store.RecordConnect(ip, clientID)
	if err != nil {
		log.Error("Failed to record visit", "error", err, "ip", ip)
		return 0
	}

	log.Info("Visitor connected",
		"ip", ip,
		"client", clientID,
		"active", t.activeCount.Load(),
	)

	return visitID
}

func (t *Tracker) OnDisconnect(visitID int64) {
	t.activeCount.Add(-1)

	if visitID > 0 {
		if err := t.store.RecordDisconnect(visitID); err != nil {
			log.Error("Failed to record disconnect", "error", err, "visitID", visitID)
		}
	}

	log.Info("Visitor disconnected", "active", t.activeCount.Load())
}

func (t *Tracker) ActiveCount() int {
	return int(t.activeCount.Load())
}

func (t *Tracker) GetStats() (Stats, error) {
	return t.store.GetStats()
}

func (t *Tracker) GetRecentVisits(limit int) ([]RecentVisit, error) {
	return t.store.GetRecentVisits(limit)
}
