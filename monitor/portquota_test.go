package monitor

import (
	"testing"
	"time"
)

func newTestQuota() *PortQuota {
	return NewPortQuota(3, 5*time.Second)
}

func TestPortQuotaNotExceededInitially(t *testing.T) {
	q := newTestQuota()
	if q.Observe(8080) {
		t.Fatal("expected first observation to not exceed quota")
	}
}

func TestPortQuotaExceededAtThreshold(t *testing.T) {
	q := newTestQuota()
	for i := 0; i < 3; i++ {
		if q.Observe(9000) {
			t.Fatalf("expected observation %d to not exceed quota", i+1)
		}
	}
	if !q.Observe(9000) {
		t.Fatal("expected 4th observation to exceed quota")
	}
}

func TestPortQuotaCountReflectsObservations(t *testing.T) {
	q := newTestQuota()
	q.Observe(443)
	q.Observe(443)
	if got := q.Count(443); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestPortQuotaIndependentPorts(t *testing.T) {
	q := newTestQuota()
	for i := 0; i < 4; i++ {
		q.Observe(80)
	}
	// port 443 should be unaffected
	if q.Observe(443) {
		t.Fatal("expected port 443 to be independent of port 80 quota")
	}
}

func TestPortQuotaResetClearsCount(t *testing.T) {
	q := newTestQuota()
	for i := 0; i < 4; i++ {
		q.Observe(8080)
	}
	q.Reset(8080)
	if q.Count(8080) != 0 {
		t.Fatal("expected count 0 after reset")
	}
	if q.Observe(8080) {
		t.Fatal("expected first observation after reset to not exceed quota")
	}
}

func TestPortQuotaEvictRemovesExpired(t *testing.T) {
	q := NewPortQuota(3, 10*time.Millisecond)
	q.Observe(22)
	time.Sleep(20 * time.Millisecond)
	q.Evict()
	if q.Count(22) != 0 {
		t.Fatal("expected expired port to be evicted")
	}
}

func TestPortQuotaSummaryEmpty(t *testing.T) {
	q := newTestQuota()
	s := q.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
	if s != "portquota: no active ports" {
		t.Fatalf("unexpected summary for empty quota: %q", s)
	}
}

func TestPortQuotaSummaryContainsPort(t *testing.T) {
	q := newTestQuota()
	q.Observe(3306)
	s := q.Summary()
	if s == "portquota: no active ports" {
		t.Fatal("expected summary to contain active port")
	}
}
