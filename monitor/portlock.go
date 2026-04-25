package monitor

import (
	"sync"
	"time"
)

// PortLockEntry holds the lock state for a single port.
type PortLockEntry struct {
	LockedAt  time.Time
	LockedBy  string
	ExpiresAt time.Time
}

// PortLock tracks ports that have been administratively locked,
// preventing any alerts or actions from being taken on them.
type PortLock struct {
	mu      sync.RWMutex
	entries map[int]PortLockEntry
	clock   func() time.Time
}

// NewPortLock creates a new PortLock.
func NewPortLock() *PortLock {
	return &PortLock{
		entries: make(map[int]PortLockEntry),
		clock:   time.Now,
	}
}

// Lock locks a port for a given duration with an optional reason string.
func (pl *PortLock) Lock(port int, by string, duration time.Duration) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	now := pl.clock()
	pl.entries[port] = PortLockEntry{
		LockedAt:  now,
		LockedBy:  by,
		ExpiresAt: now.Add(duration),
	}
}

// Unlock removes a lock from a port immediately.
func (pl *PortLock) Unlock(port int) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	delete(pl.entries, port)
}

// IsLocked returns true if the port is currently locked and the lock has not expired.
func (pl *PortLock) IsLocked(port int) bool {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	e, ok := pl.entries[port]
	if !ok {
		return false
	}
	return pl.clock().Before(e.ExpiresAt)
}

// Get returns the lock entry and whether it exists and is active.
func (pl *PortLock) Get(port int) (PortLockEntry, bool) {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	e, ok := pl.entries[port]
	if !ok || !pl.clock().Before(e.ExpiresAt) {
		return PortLockEntry{}, false
	}
	return e, true
}

// Evict removes all expired lock entries.
func (pl *PortLock) Evict() {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	now := pl.clock()
	for port, e := range pl.entries {
		if !now.Before(e.ExpiresAt) {
			delete(pl.entries, port)
		}
	}
}

// Len returns the number of active (non-expired) lock entries.
func (pl *PortLock) Len() int {
	pl.mu.RLock()
	defer pl.mu.RUnlock()
	now := pl.clock()
	count := 0
	for _, e := range pl.entries {
		if now.Before(e.ExpiresAt) {
			count++
		}
	}
	return count
}
