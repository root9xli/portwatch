package monitor

import (
	"fmt"
	"sync"
	"time"
)

// PortQuota tracks how many times a port has been seen within a rolling window
// and flags ports that exceed a configured observation quota.
type PortQuota struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	buckets map[int][]time.Time
}

// NewPortQuota creates a PortQuota that allows at most max observations of a
// single port within the given rolling window before flagging it.
func NewPortQuota(max int, window time.Duration) *PortQuota {
	return &PortQuota{
		window:  window,
		max:     max,
		buckets: make(map[int][]time.Time),
	}
}

// Observe records a new observation for port and returns true if the port has
// exceeded its quota within the rolling window.
func (q *PortQuota) Observe(port int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-q.window)

	times := q.buckets[port]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	q.buckets[port] = filtered

	return len(filtered) > q.max
}

// Count returns the number of observations for port within the current window.
func (q *PortQuota) Count(port int) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	cutoff := time.Now().Add(-q.window)
	count := 0
	for _, t := range q.buckets[port] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all observations for port.
func (q *PortQuota) Reset(port int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.buckets, port)
}

// Evict removes all ports whose last observation is older than the window.
func (q *PortQuota) Evict() {
	q.mu.Lock()
	defer q.mu.Unlock()

	cutoff := time.Now().Add(-q.window)
	for port, times := range q.buckets {
		active := false
		for _, t := range times {
			if t.After(cutoff) {
				active = true
				break
			}
		}
		if !active {
			delete(q.buckets, port)
		}
	}
}

// Summary returns a human-readable string of current quota counts.
func (q *PortQuota) Summary() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.buckets) == 0 {
		return "portquota: no active ports"
	}
	out := "portquota:\n"
	cutoff := time.Now().Add(-q.window)
	for port, times := range q.buckets {
		count := 0
		for _, t := range times {
			if t.After(cutoff) {
				count++
			}
		}
		out += fmt.Sprintf("  port %d: %d/%d\n", port, count, q.max)
	}
	return out
}
