package monitor

import (
	"sync"
	"time"
)

// PortPauseEntry holds pause state for a single port.
type PortPauseEntry struct {
	PausedAt  time.Time
	ExpiresAt time.Time
	Reason    string
}

// PortPause tracks ports that are temporarily paused from alerting.
type PortPause struct {
	mu      sync.Mutex
	entries map[int]*PortPauseEntry
	now     func() time.Time
}

// NewPortPause creates a new PortPause instance.
func NewPortPause() *PortPause {
	return &PortPause{
		entries: make(map[int]*PortPauseEntry),
		now:     time.Now,
	}
}

// Pause marks a port as paused for the given duration with an optional reason.
func (p *PortPause) Pause(port int, duration time.Duration, reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	p.entries[port] = &PortPauseEntry{
		PausedAt:  now,
		ExpiresAt: now.Add(duration),
		Reason:    reason,
	}
}

// Resume removes a pause from a port immediately.
func (p *PortPause) Resume(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.entries, port)
}

// IsPaused returns true if the port is currently paused.
func (p *PortPause) IsPaused(port int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	e, ok := p.entries[port]
	if !ok {
		return false
	}
	if p.now().After(e.ExpiresAt) {
		delete(p.entries, port)
		return false
	}
	return true
}

// Get returns the pause entry for a port, or nil if not paused.
func (p *PortPause) Get(port int) *PortPauseEntry {
	p.mu.Lock()
	defer p.mu.Unlock()
	e, ok := p.entries[port]
	if !ok {
		return nil
	}
	if p.now().After(e.ExpiresAt) {
		delete(p.entries, port)
		return nil
	}
	return e
}

// Evict removes all expired pause entries.
func (p *PortPause) Evict() {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := p.now()
	for port, e := range p.entries {
		if now.After(e.ExpiresAt) {
			delete(p.entries, port)
		}
	}
}

// Len returns the number of currently tracked (possibly expired) entries.
func (p *PortPause) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.entries)
}

// Filter returns only messages whose ports are not currently paused.
func (p *PortPause) Filter(msgs []Message) []Message {
	out := msgs[:0:0]
	for _, m := range msgs {
		if !p.IsPaused(m.Port) {
			out = append(out, m)
		}
	}
	return out
}
