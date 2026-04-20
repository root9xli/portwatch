package monitor

import "fmt"

// wellKnownPorts maps common port numbers to their service names.
var wellKnownPorts = map[int]string{
	21:    "ftp",
	22:    "ssh",
	23:    "telnet",
	25:    "smtp",
	53:    "dns",
	80:    "http",
	110:   "pop3",
	143:   "imap",
	443:   "https",
	465:   "smtps",
	587:   "submission",
	993:   "imaps",
	995:   "pop3s",
	3306:  "mysql",
	5432:  "postgres",
	6379:  "redis",
	8080:  "http-alt",
	8443:  "https-alt",
	27017: "mongodb",
}

// PortLabeler assigns human-readable labels to port numbers.
type PortLabeler struct {
	custom map[int]string
}

// NewPortLabeler creates a PortLabeler with optional custom labels merged
// on top of the built-in well-known port map.
func NewPortLabeler(custom map[int]string) *PortLabeler {
	merged := make(map[int]string, len(wellKnownPorts)+len(custom))
	for k, v := range wellKnownPorts {
		merged[k] = v
	}
	for k, v := range custom {
		merged[k] = v
	}
	return &PortLabeler{custom: merged}
}

// Label returns a descriptive string for the given port.
// If the port is recognised it returns "<port>/<service>",
// otherwise it returns the port number as a plain string.
func (pl *PortLabeler) Label(port int) string {
	if name, ok := pl.custom[port]; ok {
		return fmt.Sprintf("%d/%s", port, name)
	}
	return fmt.Sprintf("%d", port)
}

// IsWellKnown reports whether the port has a registered label.
func (pl *PortLabeler) IsWellKnown(port int) bool {
	_, ok := pl.custom[port]
	return ok
}
