package monitor

import (
	"sync"
	"time"
)

// FreezeEntry holds the freeze state for a single port.
type FreezeEntry struct {
	FrozenAt  time.Time
	ExpiresAt time.Time
	Reason    string
}

// PortFreeze tracks ports that have been "frozen" — suppressing any change
// alerts for those ports until the freeze expires or is lifted manually.
type PortFreeze struct {
	mu      sync.RWMutex
	entries map[int]*FreezeEntry
	clock   func() time.Time
}

// NewPortFreeze creates a new PortFreeze instance.
func NewPortFreeze() *PortFreeze {
	return &PortFreeze{
		entries: make(map[int]*FreezeEntry),
		clock:   time.Now,
	}
}

// Freeze marks the given port as frozen for the specified duration.
func (pf *PortFreeze) Freeze(port int, duration time.Duration, reason string) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	now := pf.clock()
	pf.entries[port] = &FreezeEntry{
		FrozenAt:  now,
		ExpiresAt: now.Add(duration),
		Reason:    reason,
	}
}

// Lift removes the freeze on the given port immediately.
func (pf *PortFreeze) Lift(port int) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	delete(pf.entries, port)
}

// IsFrozen returns true if the port is currently frozen.
func (pf *PortFreeze) IsFrozen(port int) bool {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	e, ok := pf.entries[port]
	if !ok {
		return false
	}
	return pf.clock().Before(e.ExpiresAt)
}

// Get returns the FreezeEntry for a port, or nil if not frozen.
func (pf *PortFreeze) Get(port int) *FreezeEntry {
	pf.mu.RLock()
	defer pf.mu.RUnlock()
	e, ok := pf.entries[port]
	if !ok {
		return nil
	}
	if pf.clock().Before(e.ExpiresAt) {
		return e
	}
	return nil
}

// Evict removes all expired freeze entries.
func (pf *PortFreeze) Evict() {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	now := pf.clock()
	for port, e := range pf.entries {
		if !now.Before(e.ExpiresAt) {
			delete(pf.entries, port)
		}
	}
}

// Filter removes messages for frozen ports from the provided slice.
func (pf *PortFreeze) Filter(msgs []Message) []Message {
	out := msgs[:0:len(msgs)]
	for _, m := range msgs {
		if !pf.IsFrozen(m.Port) {
			out = append(out, m)
		}
	}
	return out
}
