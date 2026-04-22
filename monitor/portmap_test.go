package monitor

import (
	"strings"
	"testing"
)

func TestPortMapSetAndGet(t *testing.T) {
	pm := NewPortMap()
	pm.Set(8080, "api-gateway", "main HTTP API")
	e, ok := pm.Get(8080)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Alias != "api-gateway" {
		t.Errorf("expected alias api-gateway, got %s", e.Alias)
	}
	if e.Comment != "main HTTP API" {
		t.Errorf("unexpected comment: %s", e.Comment)
	}
}

func TestPortMapGetMissing(t *testing.T) {
	pm := NewPortMap()
	_, ok := pm.Get(9999)
	if ok {
		t.Error("expected missing entry")
	}
}

func TestPortMapRemove(t *testing.T) {
	pm := NewPortMap()
	pm.Set(443, "https", "TLS")
	pm.Remove(443)
	_, ok := pm.Get(443)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestPortMapLen(t *testing.T) {
	pm := NewPortMap()
	if pm.Len() != 0 {
		t.Errorf("expected 0, got %d", pm.Len())
	}
	pm.Set(80, "http", "")
	pm.Set(443, "https", "")
	if pm.Len() != 2 {
		t.Errorf("expected 2, got %d", pm.Len())
	}
}

func TestPortMapSummaryEmpty(t *testing.T) {
	pm := NewPortMap()
	s := pm.Summary()
	if !strings.Contains(s, "no entries") {
		t.Errorf("expected no entries message, got: %s", s)
	}
}

func TestPortMapSummaryContainsAlias(t *testing.T) {
	pm := NewPortMap()
	pm.Set(8080, "my-service", "internal")
	s := pm.Summary()
	if !strings.Contains(s, "my-service") {
		t.Errorf("expected alias in summary, got: %s", s)
	}
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port in summary, got: %s", s)
	}
}

func TestPortMapOverwrite(t *testing.T) {
	pm := NewPortMap()
	pm.Set(3000, "old", "")
	pm.Set(3000, "new", "updated")
	e, _ := pm.Get(3000)
	if e.Alias != "new" {
		t.Errorf("expected alias new, got %s", e.Alias)
	}
}
