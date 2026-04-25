package monitor

import (
	"sync"
	"time"
)

// PortCapEntry records the peak observed connection count for a port.
type PortCapEntry struct {
	Peak      int
	ObservedAt time.Time
}

// PortCap tracks the peak (maximum) concurrent listener count seen per port
// within a rolling window. It raises a flag when a new peak is reached.
type PortCap struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[int]*PortCapEntry
}

// NewPortCap creates a PortCap with the given rolling window duration.
func NewPortCap(window time.Duration) *PortCap {
	return &PortCap{
		window:  window,
		entries: make(map[int]*PortCapEntry),
	}
}

// Observe records count for port and returns true if a new peak was reached.
func (pc *PortCap) Observe(port, count int) bool {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	e, ok := pc.entries[port]
	if !ok || now.Sub(e.ObservedAt) > pc.window {
		pc.entries[port] = &PortCapEntry{Peak: count, ObservedAt: now}
		return false
	}
	if count > e.Peak {
		e.Peak = count
		e.ObservedAt = now
		return true
	}
	return false
}

// Peak returns the current recorded peak for port, and whether it exists.
func (pc *PortCap) Peak(port int) (int, bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	e, ok := pc.entries[port]
	if !ok {
		return 0, false
	}
	return e.Peak, true
}

// Evict removes entries whose window has expired.
func (pc *PortCap) Evict() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	now := time.Now()
	for port, e := range pc.entries {
		if now.Sub(e.ObservedAt) > pc.window {
			delete(pc.entries, port)
		}
	}
}

// Len returns the number of tracked ports.
func (pc *PortCap) Len() int {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	return len(pc.entries)
}
