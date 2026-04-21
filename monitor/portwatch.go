package monitor

import (
	"fmt"
	"strings"
	"time"
)

// WatchEntry records a watched port with its watch start time and optional note.
type WatchEntry struct {
	Port      int
	AddedAt   time.Time
	Note      string
	ExpiresAt time.Time // zero means no expiry
}

// PortWatcher tracks explicitly watched ports that should always be monitored.
type PortWatcher struct {
	entries map[int]*WatchEntry
	now     func() time.Time
}

// NewPortWatcher creates a new PortWatcher.
func NewPortWatcher() *PortWatcher {
	return &PortWatcher{
		entries: make(map[int]*WatchEntry),
		now:     time.Now,
	}
}

// Watch adds a port to the watch list with an optional note and TTL.
// A zero TTL means the watch never expires.
func (w *PortWatcher) Watch(port int, note string, ttl time.Duration) {
	entry := &WatchEntry{
		Port:    port,
		AddedAt: w.now(),
		Note:    note,
	}
	if ttl > 0 {
		entry.ExpiresAt = w.now().Add(ttl)
	}
	w.entries[port] = entry
}

// Unwatch removes a port from the watch list.
func (w *PortWatcher) Unwatch(port int) {
	delete(w.entries, port)
}

// IsWatched returns true if the port is currently being watched and not expired.
func (w *PortWatcher) IsWatched(port int) bool {
	e, ok := w.entries[port]
	if !ok {
		return false
	}
	if !e.ExpiresAt.IsZero() && w.now().After(e.ExpiresAt) {
		delete(w.entries, port)
		return false
	}
	return true
}

// Evict removes all expired watch entries.
func (w *PortWatcher) Evict() {
	now := w.now()
	for port, e := range w.entries {
		if !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt) {
			delete(w.entries, port)
		}
	}
}

// Len returns the number of active (non-expired) watched ports.
func (w *PortWatcher) Len() int {
	w.Evict()
	return len(w.entries)
}

// All returns all active watch entries.
func (w *PortWatcher) All() []*WatchEntry {
	w.Evict()
	out := make([]*WatchEntry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, e)
	}
	return out
}

// Summary returns a human-readable summary of watched ports.
func (w *PortWatcher) Summary() string {
	entries := w.All()
	if len(entries) == 0 {
		return "watched ports: none"
	}
	parts := make([]string, 0, len(entries))
	for _, e := range entries {
		ttl := "permanent"
		if !e.ExpiresAt.IsZero() {
			remaining := time.Until(e.ExpiresAt).Round(time.Second)
			ttl = fmt.Sprintf("expires in %s", remaining)
		}
		parts = append(parts, fmt.Sprintf("port %d (%s) [%s]", e.Port, e.Note, ttl))
	}
	return "watched ports: " + strings.Join(parts, ", ")
}
