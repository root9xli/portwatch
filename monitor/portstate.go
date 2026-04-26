package monitor

import (
	"sync"
	"time"
)

// PortStateEntry records the current and previous state of a port.
type PortStateEntry struct {
	Port      int
	Current   string
	Previous  string
	ChangedAt time.Time
	ChangeCount int
}

// PortState tracks state transitions for ports (e.g. open -> closed -> open).
type PortState struct {
	mu      sync.Mutex
	entries map[int]*PortStateEntry
	maxAge  time.Duration
	now     func() time.Time
}

// NewPortState creates a PortState with the given max age for entries.
func NewPortState(maxAge time.Duration) *PortState {
	return &PortState{
		entries: make(map[int]*PortStateEntry),
		maxAge:  maxAge,
		now:     time.Now,
	}
}

// Transition records a state change for the given port.
// Returns true if the state actually changed.
func (ps *PortState) Transition(port int, newState string) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.evict()

	e, ok := ps.entries[port]
	if !ok {
		ps.entries[port] = &PortStateEntry{
			Port:        port,
			Current:     newState,
			Previous:    "",
			ChangedAt:   ps.now(),
			ChangeCount: 1,
		}
		return true
	}
	if e.Current == newState {
		return false
	}
	e.Previous = e.Current
	e.Current = newState
	e.ChangedAt = ps.now()
	e.ChangeCount++
	return true
}

// Get returns the entry for a port, or nil if not found.
func (ps *PortState) Get(port int) *PortStateEntry {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.evict()
	e, ok := ps.entries[port]
	if !ok {
		return nil
	}
	copy := *e
	return &copy
}

// Len returns the number of tracked ports.
func (ps *PortState) Len() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.entries)
}

// evict removes entries older than maxAge. Must be called with lock held.
func (ps *PortState) evict() {
	cutoff := ps.now().Add(-ps.maxAge)
	for port, e := range ps.entries {
		if e.ChangedAt.Before(cutoff) {
			delete(ps.entries, port)
		}
	}
}
