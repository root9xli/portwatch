package monitor

import (
	"sync"
	"time"
)

// BanEntry records why and when a port was banned.
type BanEntry struct {
	Port      int
	Reason    string
	BannedAt  time.Time
	ExpiresAt time.Time
}

// Expired returns true if the ban has passed its expiry time.
func (b BanEntry) Expired(now time.Time) bool {
	return now.After(b.ExpiresAt)
}

// PortBan tracks ports that have been flagged and temporarily banned
// from generating further alerts until the ban window expires.
type PortBan struct {
	mu      sync.Mutex
	bans    map[int]BanEntry
	window  time.Duration
	clock   func() time.Time
}

// NewPortBan creates a PortBan with the given ban window duration.
func NewPortBan(window time.Duration) *PortBan {
	return &PortBan{
		bans:   make(map[int]BanEntry),
		window: window,
		clock:  time.Now,
	}
}

// Ban records a ban for the given port with an optional reason.
func (pb *PortBan) Ban(port int, reason string) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	now := pb.clock()
	pb.bans[port] = BanEntry{
		Port:      port,
		Reason:    reason,
		BannedAt:  now,
		ExpiresAt: now.Add(pb.window),
	}
}

// IsBanned returns true if the port currently has an active ban.
func (pb *PortBan) IsBanned(port int) bool {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	entry, ok := pb.bans[port]
	if !ok {
		return false
	}
	if entry.Expired(pb.clock()) {
		delete(pb.bans, port)
		return false
	}
	return true
}

// Lift removes a ban for a port regardless of expiry.
func (pb *PortBan) Lift(port int) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	delete(pb.bans, port)
}

// Evict removes all expired bans.
func (pb *PortBan) Evict() {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	now := pb.clock()
	for port, entry := range pb.bans {
		if entry.Expired(now) {
			delete(pb.bans, port)
		}
	}
}

// Len returns the number of active (non-expired) bans.
func (pb *PortBan) Len() int {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	now := pb.clock()
	count := 0
	for _, entry := range pb.bans {
		if !entry.Expired(now) {
			count++
		}
	}
	return count
}
