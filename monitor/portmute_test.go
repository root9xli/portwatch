package monitor

import (
	"testing"
	"time"
)

func newTestMute() *PortMute {
	pm := NewPortMute()
	return pm
}

func TestPortMuteNotMutedInitially(t *testing.T) {
	pm := newTestMute()
	if pm.IsMuted(8080) {
		t.Fatal("expected port 8080 to not be muted initially")
	}
}

func TestPortMuteMutedAfterMute(t *testing.T) {
	pm := newTestMute()
	pm.Mute(8080, 5*time.Minute)
	if !pm.IsMuted(8080) {
		t.Fatal("expected port 8080 to be muted")
	}
}

func TestPortMuteUnmuteRemovesMute(t *testing.T) {
	pm := newTestMute()
	pm.Mute(8080, 5*time.Minute)
	pm.Unmute(8080)
	if pm.IsMuted(8080) {
		t.Fatal("expected port 8080 to be unmuted after Unmute")
	}
}

func TestPortMuteExpiresAfterDuration(t *testing.T) {
	now := time.Now()
	pm := NewPortMute()
	pm.nowFunc = func() time.Time { return now }
	pm.Mute(9090, 1*time.Second)
	pm.nowFunc = func() time.Time { return now.Add(2 * time.Second) }
	if pm.IsMuted(9090) {
		t.Fatal("expected mute to have expired")
	}
}

func TestPortMuteIndependentPorts(t *testing.T) {
	pm := newTestMute()
	pm.Mute(80, 5*time.Minute)
	if pm.IsMuted(443) {
		t.Fatal("port 443 should not be muted")
	}
}

func TestPortMuteFilterRemovesMuted(t *testing.T) {
	pm := newTestMute()
	pm.Mute(8080, 5*time.Minute)
	msgs := []Message{
		{Port: 8080, Level: "warn"},
		{Port: 443, Level: "info"},
	}
	filtered := pm.Filter(msgs)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 message after filter, got %d", len(filtered))
	}
	if filtered[0].Port != 443 {
		t.Fatalf("expected port 443 to pass filter, got %d", filtered[0].Port)
	}
}

func TestPortMuteLenCountsActive(t *testing.T) {
	now := time.Now()
	pm := NewPortMute()
	pm.nowFunc = func() time.Time { return now }
	pm.Mute(80, 10*time.Second)
	pm.Mute(443, 10*time.Second)
	pm.Mute(8080, 1*time.Millisecond)
	pm.nowFunc = func() time.Time { return now.Add(1 * time.Second) }
	if l := pm.Len(); l != 2 {
		t.Fatalf("expected Len 2, got %d", l)
	}
}
