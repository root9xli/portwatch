package monitor

import (
	"fmt"
	"sync"
	"time"
)

// PortDrift tracks how frequently a port's state oscillates between
// appearing and disappearing. High drift may indicate a flapping service.

type driftEntry struct {
	transitions int
	lastSeen    time.Time
	firstSeen   time.Time
}

// PortDrift records port state transitions and flags ports that change
// state more than Threshold times within Window.
type PortDrift struct {
	mu        sync.Mutex
	entries   map[int]*driftEntry
	Threshold int
	Window    time.Duration
	now       func() time.Time
}

// NewPortDrift creates a PortDrift with the given flap threshold and window.
func NewPortDrift(threshold int, window time.Duration) *PortDrift {
	return &PortDrift{
		entries:   make(map[int]*driftEntry),
		Threshold: threshold,
		Window:    window,
		now:       time.Now,
	}
}

// Observe records a state transition for the given port.
// Returns true if the port is considered flapping.
func (d *PortDrift) Observe(port int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	e, ok := d.entries[port]
	if !ok || now.Sub(e.firstSeen) > d.Window {
		d.entries[port] = &driftEntry{
			transitions: 1,
			firstSeen:   now,
			lastSeen:    now,
		}
		return false
	}
	e.transitions++
	e.lastSeen = now
	return e.transitions >= d.Threshold
}

// IsFlapping returns true if the port has exceeded the transition threshold
// within the current window.
func (d *PortDrift) IsFlapping(port int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	e, ok := d.entries[port]
	if !ok {
		return false
	}
	if d.now().Sub(e.firstSeen) > d.Window {
		return false
	}
	return e.transitions >= d.Threshold
}

// Evict removes entries whose window has expired.
func (d *PortDrift) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	for port, e := range d.entries {
		if now.Sub(e.firstSeen) > d.Window {
			delete(d.entries, port)
		}
	}
}

// Summary returns a human-readable description of flapping ports.
func (d *PortDrift) Summary() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	out := ""
	for port, e := range d.entries {
		if now.Sub(e.firstSeen) <= d.Window && e.transitions >= d.Threshold {
			out += fmt.Sprintf("port %d: %d transitions in %s\n", port, e.transitions, d.Window)
		}
	}
	if out == "" {
		return "no flapping ports"
	}
	return out
}
