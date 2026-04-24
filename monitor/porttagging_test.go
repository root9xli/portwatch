package monitor

import (
	"strings"
	"testing"
)

func TestPortTaggingAddAndGet(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(8080, "web", "internal")
	tags := pt.Get(8080)
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestPortTaggingNoDuplicates(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(443, "tls")
	pt.Tag(443, "tls", "web")
	tags := pt.Get(443)
	if len(tags) != 2 {
		t.Fatalf("expected 2 unique tags, got %d: %v", len(tags), tags)
	}
}

func TestPortTaggingHasTag(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(22, "ssh")
	if !pt.HasTag(22, "ssh") {
		t.Error("expected HasTag to return true for 'ssh'")
	}
	if pt.HasTag(22, "web") {
		t.Error("expected HasTag to return false for 'web'")
	}
}

func TestPortTaggingUntag(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(3306, "db", "mysql")
	pt.Untag(3306, "db")
	if pt.HasTag(3306, "db") {
		t.Error("expected 'db' tag to be removed")
	}
	if !pt.HasTag(3306, "mysql") {
		t.Error("expected 'mysql' tag to remain")
	}
}

func TestPortTaggingUntagAllRemovesEntry(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(9000, "debug")
	pt.Untag(9000, "debug")
	tags := pt.Get(9000)
	if len(tags) != 0 {
		t.Errorf("expected empty tags after removing all, got %v", tags)
	}
}

func TestPortTaggingGetMissingPort(t *testing.T) {
	pt := NewPortTagging()
	tags := pt.Get(1234)
	if len(tags) != 0 {
		t.Errorf("expected no tags for unknown port, got %v", tags)
	}
}

func TestPortTaggingSummaryEmpty(t *testing.T) {
	pt := NewPortTagging()
	s := pt.Summary()
	if !strings.Contains(s, "none") {
		t.Errorf("expected 'none' in empty summary, got: %s", s)
	}
}

func TestPortTaggingSummaryContainsPort(t *testing.T) {
	pt := NewPortTagging()
	pt.Tag(8080, "web")
	s := pt.Summary()
	if !strings.Contains(s, ":8080") {
		t.Errorf("expected ':8080' in summary, got: %s", s)
	}
	if !strings.Contains(s, "web") {
		t.Errorf("expected 'web' in summary, got: %s", s)
	}
}
