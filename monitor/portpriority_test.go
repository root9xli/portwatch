package monitor

import (
	"strings"
	"testing"
)

func newTestPriority() *PortPriority {
	return NewPortPriority()
}

func TestPortPrioritySetAndGet(t *testing.T) {
	pp := newTestPriority()
	pp.Set(8080, PriorityHigh, "web traffic")
	e, ok := pp.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Priority != PriorityHigh {
		t.Errorf("expected %d, got %d", PriorityHigh, e.Priority)
	}
	if e.Reason != "web traffic" {
		t.Errorf("unexpected reason: %q", e.Reason)
	}
}

func TestPortPriorityGetMissing(t *testing.T) {
	pp := newTestPriority()
	_, ok := pp.Get(9999)
	if ok {
		t.Error("expected missing entry to return false")
	}
}

func TestPortPriorityRemove(t *testing.T) {
	pp := newTestPriority()
	pp.Set(443, PriorityCritical, "tls")
	pp.Remove(443)
	_, ok := pp.Get(443)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestPortPriorityLen(t *testing.T) {
	pp := newTestPriority()
	if pp.Len() != 0 {
		t.Fatalf("expected 0, got %d", pp.Len())
	}
	pp.Set(80, PriorityNormal, "")
	pp.Set(443, PriorityCritical, "")
	if pp.Len() != 2 {
		t.Errorf("expected 2, got %d", pp.Len())
	}
}

func TestPortPriorityTopNSortedDescending(t *testing.T) {
	pp := newTestPriority()
	pp.Set(80, PriorityLow, "")
	pp.Set(443, PriorityCritical, "")
	pp.Set(8080, PriorityHigh, "")

	top := pp.TopN(0)
	if len(top) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top))
	}
	if top[0].Port != 443 {
		t.Errorf("expected port 443 first, got %d", top[0].Port)
	}
	if top[2].Port != 80 {
		t.Errorf("expected port 80 last, got %d", top[2].Port)
	}
}

func TestPortPriorityTopNLimit(t *testing.T) {
	pp := newTestPriority()
	pp.Set(80, PriorityLow, "")
	pp.Set(443, PriorityCritical, "")
	pp.Set(8080, PriorityHigh, "")

	top := pp.TopN(2)
	if len(top) != 2 {
		t.Errorf("expected 2 entries, got %d", len(top))
	}
}

func TestPortPrioritySummaryEmpty(t *testing.T) {
	pp := newTestPriority()
	s := pp.Summary()
	if !strings.Contains(s, "none") {
		t.Errorf("expected 'none' in empty summary, got %q", s)
	}
}

func TestPortPrioritySummaryContainsPort(t *testing.T) {
	pp := newTestPriority()
	pp.Set(8443, PriorityHigh, "secure api")
	s := pp.Summary()
	if !strings.Contains(s, "8443") {
		t.Errorf("expected port 8443 in summary, got %q", s)
	}
	if !strings.Contains(s, "secure api") {
		t.Errorf("expected reason in summary, got %q", s)
	}
}
