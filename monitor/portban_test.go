package monitor

import (
	"testing"
	"time"
)

func newTestBan(window time.Duration) *PortBan {
	pb := NewPortBan(window)
	return pb
}

func TestPortBanNotBannedInitially(t *testing.T) {
	pb := newTestBan(time.Minute)
	if pb.IsBanned(8080) {
		t.Fatal("expected port 8080 not to be banned initially")
	}
}

func TestPortBanBannedAfterBan(t *testing.T) {
	pb := newTestBan(time.Minute)
	pb.Ban(8080, "test reason")
	if !pb.IsBanned(8080) {
		t.Fatal("expected port 8080 to be banned")
	}
}

func TestPortBanLiftRemovesBan(t *testing.T) {
	pb := newTestBan(time.Minute)
	pb.Ban(9090, "lift test")
	pb.Lift(9090)
	if pb.IsBanned(9090) {
		t.Fatal("expected ban to be lifted")
	}
}

func TestPortBanExpiresAfterWindow(t *testing.T) {
	pb := newTestBan(time.Millisecond)
	fixed := time.Now()
	pb.clock = func() time.Time { return fixed }
	pb.Ban(7070, "expiry test")
	// advance clock past window
	pb.clock = func() time.Time { return fixed.Add(2 * time.Millisecond) }
	if pb.IsBanned(7070) {
		t.Fatal("expected ban to have expired")
	}
}

func TestPortBanIndependentPorts(t *testing.T) {
	pb := newTestBan(time.Minute)
	pb.Ban(1111, "a")
	if pb.IsBanned(2222) {
		t.Fatal("port 2222 should not be banned")
	}
	if !pb.IsBanned(1111) {
		t.Fatal("port 1111 should be banned")
	}
}

func TestPortBanEvictRemovesExpired(t *testing.T) {
	pb := newTestBan(time.Millisecond)
	fixed := time.Now()
	pb.clock = func() time.Time { return fixed }
	pb.Ban(3000, "evict test")
	pb.Ban(4000, "evict test 2")
	pb.clock = func() time.Time { return fixed.Add(5 * time.Millisecond) }
	pb.Evict()
	if pb.Len() != 0 {
		t.Fatalf("expected 0 bans after evict, got %d", pb.Len())
	}
}

func TestPortBanLen(t *testing.T) {
	pb := newTestBan(time.Minute)
	pb.Ban(100, "a")
	pb.Ban(200, "b")
	pb.Ban(300, "c")
	if pb.Len() != 3 {
		t.Fatalf("expected len 3, got %d", pb.Len())
	}
}
