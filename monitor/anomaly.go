package monitor

import (
	"sync"
	"time"
)

// AnomalyDetector tracks port churn and flags ports that appear and disappear
// repeatedly within a short observation window.
type AnomalyDetector struct {
	mu      sync.Mutex
	events  map[int][]time.Time
	window  time.Duration
	thresh  int
	now     func() time.Time
}

// NewAnomalyDetector creates a detector that flags a port as anomalous when it
// has been seen appearing/disappearing more than thresh times within window.
func NewAnomalyDetector(window time.Duration, thresh int) *AnomalyDetector {
	return &AnomalyDetector{
		events: make(map[int][]time.Time),
		window: window,
		thresh: thresh,
		now:    time.Now,
	}
}

// Record registers a churn event (add or remove) for the given port.
// It returns true if the port is now considered anomalous.
func (a *AnomalyDetector) Record(port int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := a.now()
	cutoff := now.Add(-a.window)

	evs := a.events[port]
	// prune stale events
	filtered := evs[:0]
	for _, t := range evs {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	a.events[port] = filtered

	return len(filtered) >= a.thresh
}

// IsAnomalous reports whether the port currently exceeds the churn threshold.
func (a *AnomalyDetector) IsAnomalous(port int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := a.now()
	cutoff := now.Add(-a.window)
	count := 0
	for _, t := range a.events[port] {
		if t.After(cutoff) {
			count++
		}
	}
	return count >= a.thresh
}

// Evict removes stale event records for all ports to keep memory bounded.
func (a *AnomalyDetector) Evict() {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := a.now()
	cutoff := now.Add(-a.window)
	for port, evs := range a.events {
		filtered := evs[:0]
		for _, t := range evs {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) == 0 {
			delete(a.events, port)
		} else {
			a.events[port] = filtered
		}
	}
}
