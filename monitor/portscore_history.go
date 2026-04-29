package monitor

import (
	"fmt"
	"sync"
	"time"
)

// PortScoreEntry records a scored observation for a port at a point in time.
type PortScoreEntry struct {
	Port      int
	Score     float64
	Level     string
	Timestamp time.Time
}

// PortScoreHistory maintains a rolling history of port score entries per port.
type PortScoreHistory struct {
	mu      sync.Mutex
	entries map[int][]PortScoreEntry
	maxAge  time.Duration
	maxPer  int
}

// NewPortScoreHistory creates a PortScoreHistory that retains up to maxPer
// entries per port and evicts entries older than maxAge.
func NewPortScoreHistory(maxPer int, maxAge time.Duration) *PortScoreHistory {
	return &PortScoreHistory{
		entries: make(map[int][]PortScoreEntry),
		maxAge:  maxAge,
		maxPer:  maxPer,
	}
}

// Record appends a new score entry for the given port, evicting stale entries first.
func (h *PortScoreHistory) Record(port int, score float64, level string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.evict(port)

	e := PortScoreEntry{
		Port:      port,
		Score:     score,
		Level:     level,
		Timestamp: time.Now(),
	}
	h.entries[port] = append(h.entries[port], e)

	if len(h.entries[port]) > h.maxPer {
		h.entries[port] = h.entries[port][len(h.entries[port])-h.maxPer:]
	}
}

// Get returns the recorded entries for a port, oldest first.
func (h *PortScoreHistory) Get(port int) []PortScoreEntry {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.evict(port)
	out := make([]PortScoreEntry, len(h.entries[port]))
	copy(out, h.entries[port])
	return out
}

// Summary returns a human-readable summary line for a port's score history.
func (h *PortScoreHistory) Summary(port int) string {
	entries := h.Get(port)
	if len(entries) == 0 {
		return fmt.Sprintf("port %d: no score history", port)
	}
	last := entries[len(entries)-1]
	return fmt.Sprintf("port %d: %d entries, last score=%.2f level=%s at %s",
		port, len(entries), last.Score, last.Level, last.Timestamp.Format(time.RFC3339))
}

// evict removes entries older than maxAge for the given port. Must be called with mu held.
func (h *PortScoreHistory) evict(port int) {
	cutoff := time.Now().Add(-h.maxAge)
	list := h.entries[port]
	i := 0
	for i < len(list) && list[i].Timestamp.Before(cutoff) {
		i++
	}
	if i > 0 {
		h.entries[port] = list[i:]
	}
	if len(h.entries[port]) == 0 {
		delete(h.entries, port)
	}
}
