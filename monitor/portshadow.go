package monitor

import (
	"fmt"
	"sync"
	"time"
)

// PortShadow tracks ports that appear on unexpected interfaces (e.g. a service
// that was previously only on loopback now binding on 0.0.0.0).
type PortShadow struct {
	mu      sync.Mutex
	entries map[int]shadowEntry
	maxAge  time.Duration
}

type shadowEntry struct {
	knownAddr string
	newAddr   string
	first     time.Time
	count     int
}

// ShadowEvent describes a port that has changed its bind address.
type ShadowEvent struct {
	Port      int
	KnownAddr string
	NewAddr   string
	Count     int
	First     time.Time
}

// NewPortShadow creates a PortShadow tracker with the given max age for entries.
func NewPortShadow(maxAge time.Duration) *PortShadow {
	return &PortShadow{
		entries: make(map[int]shadowEntry),
		maxAge:  maxAge,
	}
}

// Observe records that port is now seen on newAddr when it was previously on
// knownAddr. Returns a ShadowEvent and true when a shadow is detected.
func (ps *PortShadow) Observe(port int, knownAddr, newAddr string) (ShadowEvent, bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if knownAddr == newAddr {
		return ShadowEvent{}, false
	}

	e, exists := ps.entries[port]
	if !exists {
		e = shadowEntry{knownAddr: knownAddr, newAddr: newAddr, first: time.Now(), count: 1}
	} else {
		e.count++
		e.newAddr = newAddr
	}
	ps.entries[port] = e

	return ShadowEvent{
		Port:      port,
		KnownAddr: e.knownAddr,
		NewAddr:   e.newAddr,
		Count:     e.count,
		First:     e.first,
	}, true
}

// Evict removes entries older than maxAge.
func (ps *PortShadow) Evict() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	cutoff := time.Now().Add(-ps.maxAge)
	for port, e := range ps.entries {
		if e.first.Before(cutoff) {
			delete(ps.entries, port)
		}
	}
}

// Len returns the number of tracked shadow entries.
func (ps *PortShadow) Len() int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return len(ps.entries)
}

// Summary returns a human-readable report of all shadow events.
func (ps *PortShadow) Summary() string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if len(ps.entries) == 0 {
		return "portshadow: no shadow events detected"
	}
	out := fmt.Sprintf("portshadow: %d shadowed port(s)\n", len(ps.entries))
	for port, e := range ps.entries {
		out += fmt.Sprintf("  port=%d known=%s new=%s count=%d first=%s\n",
			port, e.knownAddr, e.newAddr, e.count, e.first.Format(time.RFC3339))
	}
	return out
}
