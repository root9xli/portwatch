package monitor

import (
	"sync"
	"time"
)

// DeadPort tracks ports that have been seen, gone, and reappeared
// within a short window — potential port-cycling behavior.
type DeadPort struct {
	mu      sync.Mutex
	seen    map[uint16]time.Time
	window  time.Duration
	clock   func() time.Time
}

func NewDeadPort(window time.Duration) *DeadPort {
	return &DeadPort{
		seen:   make(map[uint16]time.Time),
		window: window,
		clock:  time.Now,
	}
}

// MarkGone records when a port was last seen disappearing.
func (d *DeadPort) MarkGone(port uint16) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen[port] = d.clock()
}

// IsRecycled returns true if the port disappeared and reappeared within the window.
func (d *DeadPort) IsRecycled(port uint16) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	t, ok := d.seen[port]
	if !ok {
		return false
	}
	return d.clock().Sub(t) <= d.window
}

// Evict removes expired entries.
func (d *DeadPort) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.clock()
	for port, t := range d.seen {
		if now.Sub(t) > d.window {
			delete(d.seen, port)
		}
	}
}
