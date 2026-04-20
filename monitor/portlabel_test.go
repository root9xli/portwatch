package monitor

import "testing"

func TestPortLabelWellKnown(t *testing.T) {
	pl := NewPortLabeler(nil)
	got := pl.Label(80)
	want := "80/http"
	if got != want {
		t.Errorf("Label(80) = %q, want %q", got, want)
	}
}

func TestPortLabelUnknown(t *testing.T) {
	pl := NewPortLabeler(nil)
	got := pl.Label(9999)
	want := "9999"
	if got != want {
		t.Errorf("Label(9999) = %q, want %q", got, want)
	}
}

func TestPortLabelCustomOverride(t *testing.T) {
	pl := NewPortLabeler(map[int]string{8080: "myapp"})
	got := pl.Label(8080)
	want := "8080/myapp"
	if got != want {
		t.Errorf("Label(8080) = %q, want %q", got, want)
	}
}

func TestPortLabelCustomAddsNew(t *testing.T) {
	pl := NewPortLabeler(map[int]string{9000: "custom-svc"})
	got := pl.Label(9000)
	want := "9000/custom-svc"
	if got != want {
		t.Errorf("Label(9000) = %q, want %q", got, want)
	}
}

func TestPortLabelBuiltinPreservedWithCustom(t *testing.T) {
	pl := NewPortLabeler(map[int]string{9000: "custom-svc"})
	got := pl.Label(22)
	want := "22/ssh"
	if got != want {
		t.Errorf("Label(22) = %q, want %q", got, want)
	}
}

func TestIsWellKnownTrue(t *testing.T) {
	pl := NewPortLabeler(nil)
	if !pl.IsWellKnown(443) {
		t.Error("expected 443 to be well-known")
	}
}

func TestIsWellKnownFalse(t *testing.T) {
	pl := NewPortLabeler(nil)
	if pl.IsWellKnown(65000) {
		t.Error("expected 65000 to not be well-known")
	}
}

func TestIsWellKnownCustomPort(t *testing.T) {
	pl := NewPortLabeler(map[int]string{12345: "my-daemon"})
	if !pl.IsWellKnown(12345) {
		t.Error("expected custom port 12345 to be well-known after registration")
	}
}
