package monitor

import (
	"strings"
	"testing"
	"time"
)

func newTestShadow() *PortShadow {
	return NewPortShadow(5 * time.Minute)
}

func TestPortShadowNoEventSameAddr(t *testing.T) {
	ps := newTestShadow()
	_, ok := ps.Observe(8080, "127.0.0.1", "127.0.0.1")
	if ok {
		t.Fatal("expected no shadow event when address unchanged")
	}
}

func TestPortShadowEventOnAddrChange(t *testing.T) {
	ps := newTestShadow()
	ev, ok := ps.Observe(8080, "127.0.0.1", "0.0.0.0")
	if !ok {
		t.Fatal("expected shadow event")
	}
	if ev.Port != 8080 {
		t.Errorf("expected port 8080, got %d", ev.Port)
	}
	if ev.KnownAddr != "127.0.0.1" {
		t.Errorf("expected known addr 127.0.0.1, got %s", ev.KnownAddr)
	}
	if ev.NewAddr != "0.0.0.0" {
		t.Errorf("expected new addr 0.0.0.0, got %s", ev.NewAddr)
	}
	if ev.Count != 1 {
		t.Errorf("expected count 1, got %d", ev.Count)
	}
}

func TestPortShadowCountIncrements(t *testing.T) {
	ps := newTestShadow()
	ps.Observe(443, "127.0.0.1", "0.0.0.0")
	ev, ok := ps.Observe(443, "127.0.0.1", "0.0.0.0")
	if !ok {
		t.Fatal("expected shadow event on second observe")
	}
	if ev.Count != 2 {
		t.Errorf("expected count 2, got %d", ev.Count)
	}
}

func TestPortShadowIndependentPorts(t *testing.T) {
	ps := newTestShadow()
	ps.Observe(80, "127.0.0.1", "0.0.0.0")
	ps.Observe(443, "127.0.0.1", "0.0.0.0")
	if ps.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", ps.Len())
	}
}

func TestPortShadowEvictsExpired(t *testing.T) {
	ps := NewPortShadow(1 * time.Millisecond)
	ps.Observe(9000, "127.0.0.1", "0.0.0.0")
	time.Sleep(5 * time.Millisecond)
	ps.Evict()
	if ps.Len() != 0 {
		t.Errorf("expected 0 entries after eviction, got %d", ps.Len())
	}
}

func TestPortShadowSummaryEmpty(t *testing.T) {
	ps := newTestShadow()
	s := ps.Summary()
	if !strings.Contains(s, "no shadow") {
		t.Errorf("expected empty summary, got: %s", s)
	}
}

func TestPortShadowSummaryContainsPort(t *testing.T) {
	ps := newTestShadow()
	ps.Observe(8443, "127.0.0.1", "0.0.0.0")
	s := ps.Summary()
	if !strings.Contains(s, "8443") {
		t.Errorf("expected port 8443 in summary, got: %s", s)
	}
}
