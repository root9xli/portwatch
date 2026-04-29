package monitor

import (
	"sync"
	"time"
)

// PortVelocityEvent is emitted when a port's appearance rate exceeds a threshold.
type PortVelocityEvent struct {
	Port      int
	Rate      float64 // appearances per second
	Threshold float64
	At        time.Time
}

// PortVelocity tracks how quickly ports appear and flags those exceeding a
// rate threshold within a sliding time window.
type PortVelocity struct {
	mu        sync.Mutex
	window    time.Duration
	threshold float64
	samples   map[int][]time.Time
}

// NewPortVelocity creates a PortVelocity with the given sliding window and
// rate threshold (appearances per second).
func NewPortVelocity(window time.Duration, threshold float64) *PortVelocity {
	return &PortVelocity{
		window:    window,
		threshold: threshold,
		samples:   make(map[int][]time.Time),
	}
}

// Observe records an appearance of port at now and returns a
// PortVelocityEvent if the rate exceeds the threshold, or nil otherwise.
func (pv *PortVelocity) Observe(port int, now time.Time) *PortVelocityEvent {
	pv.mu.Lock()
	defer pv.mu.Unlock()

	cutoff := now.Add(-pv.window)
	ts := pv.samples[port]
	filtered := ts[:0]
	for _, t := range ts {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	pv.samples[port] = filtered

	if len(filtered) < 2 {
		return nil
	}
	rate := float64(len(filtered)) / pv.window.Seconds()
	if rate >= pv.threshold {
		return &PortVelocityEvent{
			Port:      port,
			Rate:      rate,
			Threshold: pv.threshold,
			At:        now,
		}
	}
	return nil
}

// Evict removes samples older than the window to free memory.
func (pv *PortVelocity) Evict(now time.Time) {
	pv.mu.Lock()
	defer pv.mu.Unlock()
	cutoff := now.Add(-pv.window)
	for port, ts := range pv.samples {
		filtered := ts[:0]
		for _, t := range ts {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}
		if len(filtered) == 0 {
			delete(pv.samples, port)
		} else {
			pv.samples[port] = filtered
		}
	}
}
