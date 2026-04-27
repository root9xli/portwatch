package monitor

import (
	"sync"
	"time"
)

// PortSpikeEvent is emitted when a port's observation rate spikes above a threshold.
type PortSpikeEvent struct {
	Port      int
	Count     int
	Threshold int
	DetectedAt time.Time
}

// PortSpike detects sudden spikes in port activity within a rolling window.
type PortSpike struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	samples   map[int][]time.Time
}

// NewPortSpike creates a PortSpike detector.
// threshold is the minimum number of observations within window to trigger a spike.
func NewPortSpike(window time.Duration, threshold int) *PortSpike {
	return &PortSpike{
		window:    window,
		threshold: threshold,
		samples:   make(map[int][]time.Time),
	}
}

// Observe records an observation for port and returns a PortSpikeEvent if a spike is detected.
func (ps *PortSpike) Observe(port int) *PortSpikeEvent {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-ps.window)

	times := ps.samples[port]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	ps.samples[port] = filtered

	if len(filtered) >= ps.threshold {
		return &PortSpikeEvent{
			Port:       port,
			Count:      len(filtered),
			Threshold:  ps.threshold,
			DetectedAt: now,
		}
	}
	return nil
}

// Reset clears all observations for port.
func (ps *PortSpike) Reset(port int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.samples, port)
}

// Len returns the number of ports currently tracked.
func (ps *PortSpike) Len() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.samples)
}
