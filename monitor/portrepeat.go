package monitor

import (
	"sync"
	"time"
)

// PortRepeatEntry tracks how many times a port has appeared in diffs.
type PortRepeatEntry struct {
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// PortRepeat counts how many times each port has been observed across diff
// cycles. It can be used to distinguish transient from persistent listeners.
type PortRepeat struct {
	mu      sync.Mutex
	entries map[int]*PortRepeatEntry
	maxAge  time.Duration
	now     func() time.Time
}

// NewPortRepeat creates a PortRepeat that evicts entries older than maxAge.
func NewPortRepeat(maxAge time.Duration) *PortRepeat {
	return &PortRepeat{
		entries: make(map[int]*PortRepeatEntry),
		maxAge:  maxAge,
		now:     time.Now,
	}
}

// Observe records an occurrence of port. Returns the updated count.
func (r *PortRepeat) Observe(port int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evict()
	e, ok := r.entries[port]
	if !ok {
		e = &PortRepeatEntry{FirstSeen: r.now()}
		r.entries[port] = e
	}
	e.Count++
	e.LastSeen = r.now()
	return e.Count
}

// Get returns the entry for port and whether it exists.
func (r *PortRepeat) Get(port int) (PortRepeatEntry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[port]
	if !ok {
		return PortRepeatEntry{}, false
	}
	return *e, true
}

// Reset clears the count for a single port.
func (r *PortRepeat) Reset(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, port)
}

// Len returns the number of tracked ports.
func (r *PortRepeat) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.entries)
}

// evict removes entries whose LastSeen is older than maxAge. Must be called
// with r.mu held.
func (r *PortRepeat) evict() {
	cutoff := r.now().Add(-r.maxAge)
	for port, e := range r.entries {
		if e.LastSeen.Before(cutoff) {
			delete(r.entries, port)
		}
	}
}
