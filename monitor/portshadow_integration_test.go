package monitor

import (
	"strings"
	"testing"
	"time"
)

func TestPortShadowDetectedOnAddrChange(t *testing.T) {
	tracker := NewPortShadow(5 * time.Minute)
	reporter := NewPortShadowReporter()

	// Simulate port 8080 moving from loopback to all-interfaces
	ev, ok := tracker.Observe(8080, "127.0.0.1", "0.0.0.0")
	if !ok {
		t.Fatal("expected shadow event")
	}
	reporter.Record(ev)

	if reporter.Len() != 1 {
		t.Errorf("expected 1 event, got %d", reporter.Len())
	}

	summary := reporter.Summary()
	if !strings.Contains(summary, "8080") {
		t.Errorf("expected port 8080 in summary: %s", summary)
	}
	if !strings.Contains(summary, "0.0.0.0") {
		t.Errorf("expected new addr in summary: %s", summary)
	}
}

func TestPortShadowTopNSortedByCount(t *testing.T) {
	tracker := NewPortShadow(5 * time.Minute)
	reporter := NewPortShadowReporter()

	// Port 443 observed 3 times, port 80 observed once
	for i := 0; i < 3; i++ {
		ev, ok := tracker.Observe(443, "127.0.0.1", "0.0.0.0")
		if ok {
			reporter.Record(ev)
		}
	}
	ev, ok := tracker.Observe(80, "127.0.0.1", "0.0.0.0")
	if ok {
		reporter.Record(ev)
	}

	top := reporter.TopN(2)
	if len(top) == 0 {
		t.Fatal("expected top results")
	}
	if top[0].Port != 443 {
		t.Errorf("expected port 443 at top, got %d", top[0].Port)
	}
}

func TestPortShadowNoEventForStableAddr(t *testing.T) {
	tracker := NewPortShadow(5 * time.Minute)
	reporter := NewPortShadowReporter()

	_, ok := tracker.Observe(3000, "0.0.0.0", "0.0.0.0")
	if ok {
		t.Fatal("should not record event when addr is unchanged")
	}

	if reporter.Len() != 0 {
		t.Errorf("expected 0 events, got %d", reporter.Len())
	}
}
