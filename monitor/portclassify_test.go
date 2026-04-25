package monitor

import "testing"

func TestPortClassifySystem(t *testing.T) {
	c := NewPortClassifier()
	if got := c.Classify(80); got != ClassSystem {
		t.Fatalf("expected system, got %s", got)
	}
}

func TestPortClassifyRegistered(t *testing.T) {
	c := NewPortClassifier()
	if got := c.Classify(8080); got != ClassRegistered {
		t.Fatalf("expected registered, got %s", got)
	}
}

func TestPortClassifyEphemeral(t *testing.T) {
	c := NewPortClassifier()
	if got := c.Classify(55000); got != ClassEphemeral {
		t.Fatalf("expected ephemeral, got %s", got)
	}
}

func TestPortClassifyOverride(t *testing.T) {
	c := NewPortClassifier()
	c.SetOverride(9090, ClassCustom)
	if got := c.Classify(9090); got != ClassCustom {
		t.Fatalf("expected custom, got %s", got)
	}
}

func TestPortClassifyRemoveOverride(t *testing.T) {
	c := NewPortClassifier()
	c.SetOverride(9090, ClassCustom)
	c.RemoveOverride(9090)
	if got := c.Classify(9090); got != ClassRegistered {
		t.Fatalf("expected registered after removal, got %s", got)
	}
}

func TestPortClassifyAll(t *testing.T) {
	c := NewPortClassifier()
	result := c.ClassifyAll([]int{22, 3000, 60000})
	if result[22] != ClassSystem {
		t.Errorf("port 22 should be system")
	}
	if result[3000] != ClassRegistered {
		t.Errorf("port 3000 should be registered")
	}
	if result[60000] != ClassEphemeral {
		t.Errorf("port 60000 should be ephemeral")
	}
}

func TestPortClassifyBoundaries(t *testing.T) {
	c := NewPortClassifier()
	cases := []struct {
		port int
		want PortClass
	}{
		{0, ClassSystem},
		{1023, ClassSystem},
		{1024, ClassRegistered},
		{49151, ClassRegistered},
		{49152, ClassEphemeral},
		{65535, ClassEphemeral},
	}
	for _, tc := range cases {
		if got := c.Classify(tc.port); got != tc.want {
			t.Errorf("port %d: want %s, got %s", tc.port, tc.want, got)
		}
	}
}
