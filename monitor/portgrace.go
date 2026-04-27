package monitor

import (
	"sync"
	"time"
)

// PortGraceEntry holds the grace period state for a port.
type PortGraceEntry struct {
	Deadline time.Time
	Port     int
}

// PortGrace suppresses alerts for newly added ports during a configurable
// grace window, allowing transient listeners to settle before alerting.
type PortGrace struct {
	mu      sync.Mutex
	entries map[int]*PortGraceEntry
	window  time.Duration
	now     func() time.Time
}

// NewPortGrace creates a PortGrace with the given grace window duration.
func NewPortGrace(window time.Duration) *PortGrace {
	return &PortGrace{
		entries: make(map[int]*PortGraceEntry),
		window:  window,
		now:     time.Now,
	}
}

// Enter registers a port as being within its grace period.
func (g *PortGrace) Enter(port int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.entries[port] = &PortGraceEntry{
		Port:     port,
		Deadline: g.now().Add(g.window),
	}
}

// InGrace reports whether the port is currently within its grace period.
func (g *PortGrace) InGrace(port int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	e, ok := g.entries[port]
	if !ok {
		return false
	}
	if g.now().After(e.Deadline) {
		delete(g.entries, port)
		return false
	}
	return true
}

// Clear removes a port from the grace registry immediately.
func (g *PortGrace) Clear(port int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.entries, port)
}

// Evict removes all expired grace entries.
func (g *PortGrace) Evict() {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.now()
	for port, e := range g.entries {
		if now.After(e.Deadline) {
			delete(g.entries, port)
		}
	}
}

// Len returns the number of ports currently in grace.
func (g *PortGrace) Len() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.entries)
}
