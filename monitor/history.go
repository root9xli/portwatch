package monitor

import (
	"sync"
	"time"
)

// HistoryEntry records a single alert event.
type HistoryEntry struct {
	Port      int
	Action    string
	Level     string
	Timestamp time.Time
}

// History keeps an in-memory ring buffer of recent alert events.
type History struct {
	mu      sync.Mutex
	entries []HistoryEntry
	maxSize int
}

// NewHistory creates a History that retains up to maxSize entries.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{maxSize: maxSize}
}

// Record appends a new entry, evicting the oldest if at capacity.
func (h *History) Record(port int, action, level string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	entry := HistoryEntry{
		Port:      port,
		Action:    action,
		Level:     level,
		Timestamp: time.Now(),
	}
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, entry)
}

// Last returns the most recent n entries (or all if fewer exist).
func (h *History) Last(n int) []HistoryEntry {
	h.mu.Lock()
	defer h.mu.Unlock()
	source := h.entries
	if n < len(h.entries) {
		source = h.entries[len(h.entries)-n:]
	}
	result := make([]HistoryEntry, len(source))
	copy(result, source)
	return result
}

// Len returns the current number of stored entries.
func (h *History) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.entries)
}

// Clear removes all entries.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = h.entries[:0]
}
