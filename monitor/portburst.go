package monitor

import (
	"sync"
	"time"
)

// BurstEvent is emitted when a port exceeds the burst threshold within a window.
type BurstEvent struct {
	Port  int
	Count int
	At    time.Time
}

// PortBurst tracks rapid successive appearances of a port within a rolling
// time window and flags it when the count exceeds a configured threshold.
type PortBurst struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	buckets   map[int][]time.Time
}

// NewPortBurst creates a PortBurst detector. window is the rolling duration;
// threshold is the minimum number of appearances within that window to
// constitute a burst.
func NewPortBurst(window time.Duration, threshold int) *PortBurst {
	if threshold < 1 {
		threshold = 1
	}
	return &PortBurst{
		window:    window,
		threshold: threshold,
		buckets:   make(map[int][]time.Time),
	}
}

// Observe records an appearance of port at now and returns a BurstEvent if
// the burst threshold has been reached, otherwise nil.
func (pb *PortBurst) Observe(port int, now time.Time) *BurstEvent {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	cutoff := now.Add(-pb.window)
	times := pb.buckets[port]

	// evict stale timestamps
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	pb.buckets[port] = filtered

	if len(filtered) >= pb.threshold {
		return &BurstEvent{Port: port, Count: len(filtered), At: now}
	}
	return nil
}

// Count returns the number of recent appearances for port within the window.
func (pb *PortBurst) Count(port int, now time.Time) int {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	cutoff := now.Add(-pb.window)
	count := 0
	for _, t := range pb.buckets[port] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Evict removes all stale entries across all ports.
func (pb *PortBurst) Evict(now time.Time) {
	pb.mu.Lock()
	defer pb.mu.Unlock()

	cutoff := now.Add(-pb.window)
	for port, times := range pb.buckets {
		filtered := times[:0]
		for _, t := range times {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) == 0 {
			delete(pb.buckets, port)
		} else {
			pb.buckets[port] = filtered
		}
	}
}
