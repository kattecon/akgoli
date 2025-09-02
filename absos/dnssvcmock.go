package absos

import (
	"net"
	"sync"
	"time"
)

// DnsMockResult represents the mock data for a DNS lookup response.
type DnsMockResult struct {
	IPs      []net.IP
	Err      error
	Duration time.Duration
}

// DnsSvcMockImpl provides a mock implementation of DnsSvc for testing purposes.
// It allows controlled DNS responses and can simulate delays using TimeSvc.
type DnsSvcMockImpl struct {
	timeSvc TimeSvc
	mu      sync.RWMutex
	results map[string]DnsMockResult
}

// NewDnsSvcMock creates a new DnsSvcMockImpl instance with the provided TimeSvc.
func NewDnsSvcMock(timeSvc TimeSvc) *DnsSvcMockImpl {
	return &DnsSvcMockImpl{
		timeSvc: timeSvc,
		results: make(map[string]DnsMockResult),
	}
}

// LookupIP performs a mock DNS lookup using predefined results.
// If a duration is set for the hostname, it will sleep using TimeSvc before returning.
func (svc *DnsSvcMockImpl) LookupIP(host string) ([]net.IP, error) {
	svc.mu.RLock()
	result, exists := svc.results[host]
	svc.mu.RUnlock()

	if !exists {
		// Return default error for unknown hosts.
		return nil, &net.DNSError{
			Err:        "no such host",
			Name:       host,
			Server:     "mock",
			IsNotFound: true,
		}
	}

	// Simulate delay if duration is set.
	if result.Duration > 0 {
		svc.timeSvc.Sleep(result.Duration)
	}

	return result.IPs, result.Err
}

// SetLookupIpResult sets the mock result for a specific hostname.
func (svc *DnsSvcMockImpl) SetLookupIpResult(host string, ips []net.IP, err error) {
	svc.SetLookupIpResultWithDuration(host, ips, err, 0)
}

// SetLookupIpResultWithDuration sets the mock result for a specific hostname with a delay duration.
func (svc *DnsSvcMockImpl) SetLookupIpResultWithDuration(host string, ips []net.IP, err error, duration time.Duration) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.results[host] = DnsMockResult{
		IPs:      ips,
		Err:      err,
		Duration: duration,
	}
}

// ClearLookupIpResult removes the mock result for a specific hostname.
func (svc *DnsSvcMockImpl) ClearLookupIpResult(host string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	delete(svc.results, host)
}

// ClearAllLookupIpResults removes all mock results.
func (svc *DnsSvcMockImpl) ClearAllLookupIpResults() {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.results = make(map[string]DnsMockResult)
}
