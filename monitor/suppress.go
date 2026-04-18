package monitor

import (
	"sync"
	"time"
)

// Suppressor prevents duplicate alerts from firing within a cooldown window.
type Suppressor struct {
	mu       sync.Mutex
	seen     map[uint16]time.Time
	cooldown time.Duration
}

// NewSuppressor creates a Suppressor with the given cooldown duration.
func NewSuppressor(cooldown time.Duration) *Suppressor {
	return &Suppressor{
		seen:     make(map[uint16]time.Time),
		cooldown: cooldown,
	}
}

// IsSuppressed returns true if an alert for the given port was already sent
// within the cooldown window. It records the port if not suppressed.
func (s *Suppressor) IsSuppressed(port uint16) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if last, ok := s.seen[port]; ok {
		if time.Since(last) < s.cooldown {
			return true
		}
	}
	s.seen[port] = time.Now()
	return false
}

// Expire removes stale entries older than the cooldown window.
func (s *Suppressor) Expire() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for port, last := range s.seen {
		if time.Since(last) >= s.cooldown {
			delete(s.seen, port)
		}
	}
}
