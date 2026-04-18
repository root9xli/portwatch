package monitor

import (
	"sync"
	"time"
)

// RateLimiter limits how many alerts can be sent per port within a window.
type RateLimiter struct {
	mu      sync.Mutex
	counts  map[int][]time.Time
	window  time.Duration
	maxHits int
}

// NewRateLimiter creates a RateLimiter that allows at most maxHits alerts
// per port within the given window duration.
func NewRateLimiter(window time.Duration, maxHits int) *RateLimiter {
	return &RateLimiter{
		counts:  make(map[int][]time.Time),
		window:  window,
		maxHits: maxHits,
	}
}

// Allow returns true if an alert for the given port is within the rate limit.
func (r *RateLimiter) Allow(port int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.window)

	times := r.counts[port]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= r.maxHits {
		r.counts[port] = filtered
		return false
	}

	r.counts[port] = append(filtered, now)
	return true
}

// Reset clears the rate limit state for a specific port.
func (r *RateLimiter) Reset(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.counts, port)
}
