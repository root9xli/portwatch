package monitor

import (
	"testing"
	"time"
)

func TestPortConnIntegrationWithDiff(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	snap1 := Snapshot{}
	snap2 := Snapshot{Ports: []PortEntry{{Port: port, State: "LISTEN"}}}

	msgs := Diff(snap1, snap2)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	pc := NewPortConn(200 * time.Millisecond)
	errs := pc.ProbeAll([]int{msgs[0].Port})
	if len(errs) != 0 {
		t.Errorf("expected port from diff to be reachable, got errors: %v", errs)
	}

	rec, ok := pc.Get(msgs[0].Port)
	if !ok {
		t.Fatal("expected conn record after probe")
	}
	if rec.Count < 1 {
		t.Errorf("expected at least 1 connection recorded, got %d", rec.Count)
	}
}

func TestPortConnIntegrationNoProbeOnRemoved(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	snap1 := Snapshot{Ports: []PortEntry{{Port: port, State: "LISTEN"}}}
	snap2 := Snapshot{}

	msgs := Diff(snap1, snap2)
	if len(msgs) != 1 {
		t.Fatalf("expected 1 diff message, got %d", len(msgs))
	}

	// Only probe added ports; removed ports should not be probed.
	var added []int
	for _, m := range msgs {
		if m.Action == "added" {
			added = append(added, m.Port)
		}
	}

	pc := NewPortConn(100 * time.Millisecond)
	if len(added) > 0 {
		pc.ProbeAll(added)
	}

	_, ok := pc.Get(port)
	if ok {
		t.Error("removed port should not have a conn record")
	}
}
