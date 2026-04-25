package monitor

import "sync"

// PortClass represents a category assigned to a port.
type PortClass string

const (
	ClassSystem    PortClass = "system"    // 0–1023
	ClassRegistered PortClass = "registered" // 1024–49151
	ClassEphemeral PortClass = "ephemeral" // 49152–65535
	ClassCustom    PortClass = "custom"    // user-defined override
)

// PortClassifier assigns a PortClass to a port number.
type PortClassifier struct {
	mu      sync.RWMutex
	override map[int]PortClass
}

// NewPortClassifier returns a PortClassifier with an empty override map.
func NewPortClassifier() *PortClassifier {
	return &PortClassifier{
		override: make(map[int]PortClass),
	}
}

// SetOverride assigns a custom class to a specific port.
func (c *PortClassifier) SetOverride(port int, class PortClass) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.override[port] = class
}

// RemoveOverride removes any custom override for the given port.
func (c *PortClassifier) RemoveOverride(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.override, port)
}

// Classify returns the PortClass for the given port number.
func (c *PortClassifier) Classify(port int) PortClass {
	c.mu.RLock()
	if cls, ok := c.override[port]; ok {
		c.mu.RUnlock()
		return cls
	}
	c.mu.RUnlock()

	switch {
	case port >= 0 && port <= 1023:
		return ClassSystem
	case port >= 1024 && port <= 49151:
		return ClassRegistered
	default:
		return ClassEphemeral
	}
}

// ClassifyAll returns a map of port → PortClass for each port in the slice.
func (c *PortClassifier) ClassifyAll(ports []int) map[int]PortClass {
	out := make(map[int]PortClass, len(ports))
	for _, p := range ports {
		out[p] = c.Classify(p)
	}
	return out
}
