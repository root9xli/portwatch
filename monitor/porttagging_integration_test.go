package monitor

import (
	"strings"
	"testing"
)

func TestPortTaggingReporterTaggedPorts(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(80, "web")
	pt.Tag(443, "tls", "web")
	pt.Tag(22, "ssh")

	r := NewPortTaggingReporter(pt)
	ports := r.TaggedPorts()

	if len(ports) != 3 {
		t.Fatalf("expected 3 tagged ports, got %d", len(ports))
	}
	// Should be sorted
	if ports[0] != 22 || ports[1] != 80 || ports[2] != 443 {
		t.Errorf("unexpected order: %v", ports)
	}
}

func TestPortTaggingReporterPortsWithTag(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(80, "web")
	pt.Tag(443, "web", "tls")
	pt.Tag(22, "ssh")

	r := NewPortTaggingReporter(pt)
	webPorts := r.PortsWithTag("web")

	if len(webPorts) != 2 {
		t.Fatalf("expected 2 ports with 'web' tag, got %d: %v", len(webPorts), webPorts)
	}
	if webPorts[0] != 80 || webPorts[1] != 443 {
		t.Errorf("unexpected web ports: %v", webPorts)
	}
}

func TestPortTaggingReporterSummaryAfterTagging(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(8080, "dev", "http")

	r := NewPortTaggingReporter(pt)
	s := r.Summary()

	if !strings.Contains(s, "8080") {
		t.Errorf("expected port 8080 in summary, got: %s", s)
	}
	if !strings.Contains(s, "dev") {
		t.Errorf("expected tag 'dev' in summary, got: %s", s)
	}
	if !strings.Contains(s, "http") {
		t.Errorf("expected tag 'http' in summary, got: %s", s)
	}
}

func TestPortTaggingReporterEmptySummary(t *testing.T) {
	pt := NewPortTagging()
	r := NewPortTaggingReporter(pt)
	s := r.Summary()
	if !strings.Contains(s, "no tagged ports") {
		t.Errorf("expected empty message, got: %s", s)
	}
}

func TestPortTaggingReporterUntagReflectedInReport(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(5432, "db", "postgres")
	pt.Untag(5432, "db")

	r := NewPortTaggingReporter(pt)
	dbPorts := r.PortsWithTag("db")
	if len(dbPorts) != 0 {
		t.Errorf("expected no ports with 'db' after untag, got: %v", dbPorts)
	}
	postgresPorts := r.PortsWithTag("postgres")
	if len(postgresPorts) != 1 || postgresPorts[0] != 5432 {
		t.Errorf("expected port 5432 with 'postgres' tag, got: %v", postgresPorts)
	}
}
