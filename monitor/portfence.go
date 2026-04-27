package monitor

import (
	"fmt"
	"sync"
	"time"
)

// PortFenceEntry records when a port was fenced and for how long.
type PortFenceEntry struct {
	Port      int
	FencedAt  time.Time
	ExpiresAt time.Time
	Reason    string
}

// PortFence tracks ports that are temporarily fenced (blocked from generating
// new alerts) for a configurable duration, with an optional reason.
type PortFence struct {
	mu      sync.Mutex
	entries map[int]*PortFenceEntry
	clock   func() time.Time
}

// NewPortFence creates a new PortFence using real wall-clock time.
func NewPortFence() *PortFence {
	return &PortFence{
		entries: make(map[int]*PortFenceEntry),
		clock:   time.Now,
	}
}

// Fence marks the given port as fenced for the specified duration.
func (f *PortFence) Fence(port int, duration time.Duration, reason string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	f.entries[port] = &PortFenceEntry{
		Port:      port,
		FencedAt:  now,
		ExpiresAt: now.Add(duration),
		Reason:    reason,
	}
}

// Lift removes a fence from the given port immediately.
func (f *PortFence) Lift(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.entries, port)
}

// IsFenced returns true if the port is currently fenced.
func (f *PortFence) IsFenced(port int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	e, ok := f.entries[port]
	if !ok {
		return false
	}
	if f.clock().After(e.ExpiresAt) {
		delete(f.entries, port)
		return false
	}
	return true
}

// Get returns the fence entry for a port, or nil if not fenced.
func (f *PortFence) Get(port int) *PortFenceEntry {
	f.mu.Lock()
	defer f.mu.Unlock()
	e, ok := f.entries[port]
	if !ok {
		return nil
	}
	if f.clock().After(e.ExpiresAt) {
		delete(f.entries, port)
		return nil
	}
	return e
}

// Len returns the number of currently active (non-expired) fences.
func (f *PortFence) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	count := 0
	for port, e := range f.entries {
		if now.After(e.ExpiresAt) {
			delete(f.entries, port)
		} else {
			count++
		}
	}
	return count
}

// Filter removes messages for fenced ports and returns the rest.
func (f *PortFence) Filter(msgs []Message) []Message {
	out := msgs[:0:0]
	for _, m := range msgs {
		if !f.IsFenced(m.Port) {
			out = append(out, m)
		}
	}
	return out
}

// Summary returns a human-readable description of all active fences.
func (f *PortFence) Summary() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	now := f.clock()
	if len(f.entries) == 0 {
		return "portfence: no active fences"
	}
	s := fmt.Sprintf("portfence: %d active fence(s)\n", len(f.entries))
	for _, e := range f.entries {
		if now.After(e.ExpiresAt) {
			continue
		}
		remaining := e.ExpiresAt.Sub(now).Round(time.Second)
		s += fmt.Sprintf("  port %d fenced for %s (reason: %s)\n", e.Port, remaining, e.Reason)
	}
	return s
}
