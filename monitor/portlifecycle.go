package monitor

import (
	"sync"
	"time"
)

// PortLifecycle tracks how long each port has been continuously observed
// as open, and records when it was first and last seen.
type PortLifecycle struct {
	mu      sync.Mutex
	records map[int]*lifecycleRecord
}

type lifecycleRecord struct {
	FirstSeen time.Time
	LastSeen  time.Time
	SeenCount int
}

// LifecycleEntry is a snapshot of a single port's lifecycle data.
type LifecycleEntry struct {
	Port      int
	FirstSeen time.Time
	LastSeen  time.Time
	Uptime    time.Duration
	SeenCount int
}

// NewPortLifecycle creates a new PortLifecycle tracker.
func NewPortLifecycle() *PortLifecycle {
	return &PortLifecycle{
		records: make(map[int]*lifecycleRecord),
	}
}

// Observe records that a port was seen at the current time.
func (pl *PortLifecycle) Observe(port int) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	now := time.Now()
	if r, ok := pl.records[port]; ok {
		r.LastSeen = now
		r.SeenCount++
	} else {
		pl.records[port] = &lifecycleRecord{
			FirstSeen: now,
			LastSeen:  now,
			SeenCount: 1,
		}
	}
}

// Forget removes a port's lifecycle record (e.g. when it is no longer open).
func (pl *PortLifecycle) Forget(port int) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	delete(pl.records, port)
}

// Get returns the LifecycleEntry for a port, and whether it exists.
func (pl *PortLifecycle) Get(port int) (LifecycleEntry, bool) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	r, ok := pl.records[port]
	if !ok {
		return LifecycleEntry{}, false
	}
	return LifecycleEntry{
		Port:      port,
		FirstSeen: r.FirstSeen,
		LastSeen:  r.LastSeen,
		Uptime:    r.LastSeen.Sub(r.FirstSeen),
		SeenCount: r.SeenCount,
	}, true
}

// All returns lifecycle entries for all currently tracked ports.
func (pl *PortLifecycle) All() []LifecycleEntry {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	out := make([]LifecycleEntry, 0, len(pl.records))
	for port, r := range pl.records {
		out = append(out, LifecycleEntry{
			Port:      port,
			FirstSeen: r.FirstSeen,
			LastSeen:  r.LastSeen,
			Uptime:    r.LastSeen.Sub(r.FirstSeen),
			SeenCount: r.SeenCount,
		})
	}
	return out
}
