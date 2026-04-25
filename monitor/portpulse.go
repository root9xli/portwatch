package monitor

import (
	"sync"
	"time"
)

// PortPulse tracks the last-seen timestamp for each observed port,
// allowing callers to detect ports that have gone silent (no recent activity).
type PortPulse struct {
	mu      sync.Mutex
	pulses  map[int]time.Time
	maxAge  time.Duration
	now     func() time.Time
}

// NewPortPulse creates a PortPulse that considers a port silent after maxAge
// without an observation.
func NewPortPulse(maxAge time.Duration) *PortPulse {
	return &PortPulse{
		pulses: make(map[int]time.Time),
		maxAge: maxAge,
		now:    time.Now,
	}
}

// Observe records the current time as the last-seen timestamp for port.
func (p *PortPulse) Observe(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pulses[port] = p.now()
}

// Silent returns true if port has not been observed within maxAge,
// or has never been observed.
func (p *PortPulse) Silent(port int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	t, ok := p.pulses[port]
	if !ok {
		return true
	}
	return p.now().Sub(t) > p.maxAge
}

// LastSeen returns the last observation time for port and whether it exists.
func (p *PortPulse) LastSeen(port int) (time.Time, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	t, ok := p.pulses[port]
	return t, ok
}

// Evict removes ports whose last observation exceeds maxAge.
func (p *PortPulse) Evict() {
	p.mu.Lock()
	defer p.mu.Unlock()
	cutoff := p.now().Add(-p.maxAge)
	for port, t := range p.pulses {
		if t.Before(cutoff) {
			delete(p.pulses, port)
		}
	}
}

// Len returns the number of tracked ports.
func (p *PortPulse) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.pulses)
}
