package monitor

import (
	"sync"
	"time"
)

// PortWindowEntry holds the observation window bounds for a port.
type PortWindowEntry struct {
	Start    time.Time
	End      time.Time
	Port     int
}

// PortWindow tracks time-bounded observation windows per port.
// A port is considered "in window" if the current time falls between Start and End.
type PortWindow struct {
	mu      sync.Mutex
	windows map[int]PortWindowEntry
	now     func() time.Time
}

// NewPortWindow creates a new PortWindow.
func NewPortWindow() *PortWindow {
	return &PortWindow{
		windows: make(map[int]PortWindowEntry),
		now:     time.Now,
	}
}

// Register sets an observation window [start, end] for the given port.
func (pw *PortWindow) Register(port int, start, end time.Time) {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	pw.windows[port] = PortWindowEntry{
		Port:  port,
		Start: start,
		End:   end,
	}
}

// InWindow reports whether the given port has an active window at the current time.
func (pw *PortWindow) InWindow(port int) bool {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	e, ok := pw.windows[port]
	if !ok {
		return false
	}
	now := pw.now()
	return !now.Before(e.Start) && now.Before(e.End)
}

// Remove deletes the window entry for the given port.
func (pw *PortWindow) Remove(port int) {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	delete(pw.windows, port)
}

// Evict removes all expired window entries (where End is in the past).
func (pw *PortWindow) Evict() {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	now := pw.now()
	for port, e := range pw.windows {
		if now.After(e.End) || now.Equal(e.End) {
			delete(pw.windows, port)
		}
	}
}

// Len returns the number of registered window entries.
func (pw *PortWindow) Len() int {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	return len(pw.windows)
}

// Get returns the window entry for a port and whether it exists.
func (pw *PortWindow) Get(port int) (PortWindowEntry, bool) {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	e, ok := pw.windows[port]
	return e, ok
}
