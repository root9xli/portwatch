package monitor

import (
	"sync"
	"time"
)

// EvictedPort records when a port was last seen and when it was evicted.
type EvictedPort struct {
	Port      int
	LastSeen  time.Time
	EvictedAt time.Time
	Count     int
}

// PortEvict tracks ports that have been evicted (removed) and retains a
// rolling history so reporters can surface recently-gone listeners.
type PortEvict struct {
	mu      sync.Mutex
	records map[int]*EvictedPort
	maxAge  time.Duration
	now     func() time.Time
}

// NewPortEvict creates a PortEvict that retains eviction records for maxAge.
func NewPortEvict(maxAge time.Duration) *PortEvict {
	return &PortEvict{
		records: make(map[int]*EvictedPort),
		maxAge:  maxAge,
		now:     time.Now,
	}
}

// Evict records that port was removed at the current time.
func (pe *PortEvict) Evict(port int, lastSeen time.Time) {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.prune()
	if r, ok := pe.records[port]; ok {
		r.Count++
		r.LastSeen = lastSeen
		r.EvictedAt = pe.now()
		return
	}
	pe.records[port] = &EvictedPort{
		Port:      port,
		LastSeen:  lastSeen,
		EvictedAt: pe.now(),
		Count:     1,
	}
}

// Get returns the eviction record for port, if present and not expired.
func (pe *PortEvict) Get(port int) (*EvictedPort, bool) {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.prune()
	r, ok := pe.records[port]
	if !ok {
		return nil, false
	}
	copy := *r
	return &copy, true
}

// All returns all non-expired eviction records.
func (pe *PortEvict) All() []EvictedPort {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.prune()
	out := make([]EvictedPort, 0, len(pe.records))
	for _, r := range pe.records {
		out = append(out, *r)
	}
	return out
}

// Len returns the number of non-expired eviction records.
func (pe *PortEvict) Len() int {
	pe.mu.Lock()
	defer pe.mu.Unlock()
	pe.prune()
	return len(pe.records)
}

// prune removes records older than maxAge. Caller must hold mu.
func (pe *PortEvict) prune() {
	cutoff := pe.now().Add(-pe.maxAge)
	for k, r := range pe.records {
		if r.EvictedAt.Before(cutoff) {
			delete(pe.records, k)
		}
	}
}
