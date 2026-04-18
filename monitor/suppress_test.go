package monitor

import (
	"testing"
	"time"
)

func TestSuppressorNotSuppressedFirstTime(t *testing.T) {
	s := NewSuppressor(5 * time.Second)
	if s.IsSuppressed(8080) {
		t.Fatal("expected port 8080 to not be suppressed on first call")
	}
}

func TestSuppressorSuppressedSecondCall(t *testing.T) {
	s := NewSuppressor(5 * time.Second)
	s.IsSuppressed(9090)
	if !s.IsSuppressed(9090) {
		t.Fatal("expected port 9090 to be suppressed on second call")
	}
}

func TestSuppressorDifferentPortsIndependent(t *testing.T) {
	s := NewSuppressor(5 * time.Second)
	s.IsSuppressed(1111)
	if s.IsSuppressed(2222) {
		t.Fatal("expected port 2222 to not be suppressed")
	}
}

func TestSuppressorExpiresAfterCooldown(t *testing.T) {
	s := NewSuppressor(50 * time.Millisecond)
	s.IsSuppressed(3000)
	time.Sleep(80 * time.Millisecond)
	if s.IsSuppressed(3000) {
		t.Fatal("expected port 3000 to no longer be suppressed after cooldown")
	}
}

func TestSuppressorExpireRemovesStale(t *testing.T) {
	s := NewSuppressor(50 * time.Millisecond)
	s.IsSuppressed(4000)
	time.Sleep(80 * time.Millisecond)
	s.Expire()
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.seen[4000]; ok {
		t.Fatal("expected stale entry for port 4000 to be removed")
	}
}
