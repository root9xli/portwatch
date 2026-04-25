package monitor

import (
	"testing"
	"time"
)

func newTestResolve(ttl time.Duration) *PortResolve {
	pr := NewPortResolve(ttl)
	// Stub nowFunc so we control time in tests.
	pr.nowFunc = time.Now
	return pr
}

func TestPortResolveCachesResult(t *testing.T) {
	pr := newTestResolve(5 * time.Minute)
	// Seed cache directly to avoid real DNS.
	pr.mu.Lock()
	pr.cache[8080] = resolveEntry{
		hostnames:  []string{"localhost."},
		resolvedAt: pr.nowFunc(),
	}
	pr.mu.Unlock()

	hosts := pr.Get(8080)
	if len(hosts) == 0 {
		t.Fatal("expected cached hostnames, got none")
	}
	if hosts[0] != "localhost." {
		t.Errorf("expected localhost., got %s", hosts[0])
	}
}

func TestPortResolveGetMissingReturnsNil(t *testing.T) {
	pr := newTestResolve(5 * time.Minute)
	if got := pr.Get(9999); got != nil {
		t.Errorf("expected nil for unknown port, got %v", got)
	}
}

func TestPortResolveExpiredEntryReturnsNil(t *testing.T) {
	pr := newTestResolve(1 * time.Millisecond)
	past := time.Now().Add(-1 * time.Second)
	pr.mu.Lock()
	pr.cache[443] = resolveEntry{
		hostnames:  []string{"example.com."},
		resolvedAt: past,
	}
	pr.mu.Unlock()

	if got := pr.Get(443); got != nil {
		t.Errorf("expected nil for expired entry, got %v", got)
	}
}

func TestPortResolveEvictRemovesExpired(t *testing.T) {
	pr := newTestResolve(1 * time.Millisecond)
	past := time.Now().Add(-1 * time.Second)
	pr.mu.Lock()
	pr.cache[80] = resolveEntry{hostnames: []string{"old."}, resolvedAt: past}
	pr.cache[22] = resolveEntry{hostnames: []string{"fresh."}, resolvedAt: time.Now().Add(time.Hour)}
	pr.mu.Unlock()

	pr.Evict()

	if pr.Len() != 1 {
		t.Errorf("expected 1 entry after evict, got %d", pr.Len())
	}
	if pr.Get(80) != nil {
		t.Error("expected evicted entry to be gone")
	}
}

func TestPortResolveLenReflectsCache(t *testing.T) {
	pr := newTestResolve(5 * time.Minute)
	if pr.Len() != 0 {
		t.Errorf("expected 0 initially, got %d", pr.Len())
	}
	pr.mu.Lock()
	pr.cache[3000] = resolveEntry{hostnames: []string{}, resolvedAt: time.Now()}
	pr.mu.Unlock()
	if pr.Len() != 1 {
		t.Errorf("expected 1 after insert, got %d", pr.Len())
	}
}
