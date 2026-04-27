package monitor

import (
	"sync"
	"time"
)

// PortFlapEvent is emitted when a port is detected as flapping.
type PortFlapEvent struct {
	Port      int
	Flaps     int
	Window    time.Duration
	DetectedAt time.Time
}

// PortFlap tracks ports that repeatedly appear and disappear within a window.
type PortFlap struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    map[int][]time.Time
}

// NewPortFlap creates a PortFlap detector.
// threshold is the minimum number of transitions within window to flag a flap.
func NewPortFlap(threshold int, window time.Duration) *PortFlap {
	return &PortFlap{
		window:    window,
		threshold: threshold,
		events:    make(map[int][]time.Time),
	}
}

// Observe records a transition (add or remove) for the given port.
// Returns a PortFlapEvent if the port is now considered flapping, nil otherwise.
func (f *PortFlap) Observe(port int, now time.Time) *PortFlapEvent {
	f.mu.Lock()
	defer f.mu.Unlock()

	cutoff := now.Add(-f.window)
	times := f.events[port]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	f.events[port] = filtered

	if len(filtered) >= f.threshold {
		return &PortFlapEvent{
			Port:       port,
			Flaps:      len(filtered),
			Window:     f.window,
			DetectedAt: now,
		}
	}
	return nil
}

// Forget removes tracking data for a port.
func (f *PortFlap) Forget(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.events, port)
}

// Evict removes stale entries for all ports whose last event is outside the window.
func (f *PortFlap) Evict(now time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()
	cutoff := now.Add(-f.window)
	for port, times := range f.events {
		if len(times) == 0 || times[len(times)-1].Before(cutoff) {
			delete(f.events, port)
		}
	}
}
