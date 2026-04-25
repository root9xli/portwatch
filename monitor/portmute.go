package monitor

import (
	"sync"
	"time"
)

// PortMuteEntry holds the expiry time for a muted port.
type PortMuteEntry struct {
	Until time.Time
}

// PortMute allows temporarily muting alerts for specific ports for a
// configurable duration. Muted ports are silently skipped during alert
// dispatch. Entries are lazily evicted on access.
type PortMute struct {
	mu      sync.Mutex
	muted   map[int]PortMuteEntry
	nowFunc func() time.Time
}

// NewPortMute creates a new PortMute instance.
func NewPortMute() *PortMute {
	return &PortMute{
		muted:   make(map[int]PortMuteEntry),
		nowFunc: time.Now,
	}
}

// Mute silences alerts for the given port for the specified duration.
func (pm *PortMute) Mute(port int, duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.muted[port] = PortMuteEntry{Until: pm.nowFunc().Add(duration)}
}

// Unmute removes the mute for the given port immediately.
func (pm *PortMute) Unmute(port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.muted, port)
}

// IsMuted returns true if the port is currently muted.
func (pm *PortMute) IsMuted(port int) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	entry, ok := pm.muted[port]
	if !ok {
		return false
	}
	if pm.nowFunc().After(entry.Until) {
		delete(pm.muted, port)
		return false
	}
	return true
}

// Filter removes messages whose port is currently muted.
func (pm *PortMute) Filter(msgs []Message) []Message {
	out := msgs[:0:0]
	for _, m := range msgs {
		if !pm.IsMuted(m.Port) {
			out = append(out, m)
		}
	}
	return out
}

// Len returns the number of currently active (non-expired) mutes.
func (pm *PortMute) Len() int {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	now := pm.nowFunc()
	count := 0
	for port, entry := range pm.muted {
		if now.After(entry.Until) {
			delete(pm.muted, port)
		} else {
			count++
		}
	}
	return count
}
