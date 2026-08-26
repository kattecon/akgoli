// Package absos provides interfaces that abstract operating-system services
// such as DNS resolution and the system clock. Tests can replace the real
// implementations with mocks to control network responses and time
// progression without real I/O or delays.
package absos

import "net"

// DnsSvc abstracts DNS lookups. The production implementation delegates to
// net.LookupIP. Use NewDnsSvcMock for a controlled test replacement.
type DnsSvc interface {
	// LookupIP resolves host to a list of IP addresses. It follows the same
	// rules as net.LookupIP, including the returned error types.
	LookupIP(host string) ([]net.IP, error)
}

type dnsSvcImpl struct{}

var dnsSvcImplInstance = dnsSvcImpl{}

// NewDnsSvc returns the production DnsSvc that delegates to net.LookupIP.
// The returned value is stateless and shared across callers.
func NewDnsSvc() DnsSvc {
	return dnsSvcImplInstance
}

func (dnsSvcImpl) LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}
