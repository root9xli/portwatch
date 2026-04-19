package monitor

import (
	"sync"
	"time"
)

// Throttle limits how frequently alerts fire per port using a token bucket.
type Throttle struct {
	mu       sync.Mutex
	tokens   map[int]int
	lastFill map[int]time.Time
	max      int
	refill   time.Duration
}

// NewThrottle creates a Throttle allowing max alerts per refill window per port.
func NewThrottle(max int, refill time.Duration) *Throttle {
	return &Throttle{
		tokens:   make(map[int]int),
		lastFill: make(map[int]time.Time),
		max:      max,
		refill:   refill,
	}
}

// Allow returns true if the alert for the given port should proceed.
func (t *Throttle) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	last, seen := t.lastFill[port]
	if !seen || now.Sub(last) >= t.refill {
		t.tokens[port] = t.max
		t.lastFill[port] = now
	}

	if t.tokens[port] <= 0 {
		return false
	}
	t.tokens[port]--
	return true
}

// Reset clears throttle state for a port.
func (t *Throttle) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tokens, port)
	delete(t.lastFill, port)
}

// Purge removes state for ports not seen since the given cutoff.
func (t *Throttle) Purge(cutoff time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for port, last := range t.lastFill {
		if last.Before(cutoff) {
			delete(t.tokens, port)
			delete(t.lastFill, port)
		}
	}
}
