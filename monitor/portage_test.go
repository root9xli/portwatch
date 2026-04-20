package monitor

import (
	"strings"
	"testing"
	"time"
)

func newTestPortAge() (*PortAge, *time.Time) {
	now := time.Now()
	pa := NewPortAge()
	pa.clock = func() time.Time { return now }
	return pa, &now
}

func TestPortAgeObserveAndAge(t *testing.T) {
	pa, now := newTestPortAge()
	pa.Observe(8080)
	*now = now.Add(5 * time.Second)
	age, ok := pa.Age(8080)
	if !ok {
		t.Fatal("expected port to be known")
	}
	if age != 5*time.Second {
		t.Fatalf("expected 5s, got %s", age)
	}
}

func TestPortAgeUnknownPort(t *testing.T) {
	pa, _ := newTestPortAge()
	_, ok := pa.Age(9999)
	if ok {
		t.Fatal("expected unknown port to return ok=false")
	}
}

func TestPortAgeObserveIdempotent(t *testing.T) {
	pa, now := newTestPortAge()
	pa.Observe(443)
	*now = now.Add(10 * time.Second)
	pa.Observe(443) // second observe should not reset the clock
	age, ok := pa.Age(443)
	if !ok {
		t.Fatal("expected port to be known")
	}
	if age != 10*time.Second {
		t.Fatalf("expected 10s, got %s", age)
	}
}

func TestPortAgeForget(t *testing.T) {
	pa, _ := newTestPortAge()
	pa.Observe(22)
	pa.Forget(22)
	_, ok := pa.Age(22)
	if ok {
		t.Fatal("expected forgotten port to be unknown")
	}
}

func TestPortAgeLen(t *testing.T) {
	pa, _ := newTestPortAge()
	pa.Observe(80)
	pa.Observe(443)
	pa.Observe(8080)
	if pa.Len() != 3 {
		t.Fatalf("expected 3, got %d", pa.Len())
	}
	pa.Forget(443)
	if pa.Len() != 2 {
		t.Fatalf("expected 2 after forget, got %d", pa.Len())
	}
}

func TestPortAgeSummaryContainsPort(t *testing.T) {
	pa, now := newTestPortAge()
	pa.Observe(8080)
	*now = now.Add(2 * time.Minute)
	summary := pa.Summary()
	if !strings.Contains(summary, "8080") {
		t.Fatalf("expected summary to contain port 8080, got: %s", summary)
	}
}

func TestPortAgeSummaryEmpty(t *testing.T) {
	pa, _ := newTestPortAge()
	summary := pa.Summary()
	if summary != "no ports tracked" {
		t.Fatalf("expected empty summary message, got: %s", summary)
	}
}
