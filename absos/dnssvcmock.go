package absos

import (
	"net"
	"sync"
	"time"
)

// DnsMockResult holds the data returned by a mock DNS lookup.
type DnsMockResult struct {
	// IPs is the list of addresses returned by LookupIP. The mock stores and
	// returns this slice without copying, so the caller and the mock share
	// the same backing array. Mutating the slice after registration changes
	// future lookup results.
	IPs []net.IP

	// Err is the error returned by LookupIP alongside IPs.
	Err error

	// Duration is the delay applied before returning the result. A positive
	// value causes LookupIP to sleep using the TimeSvc passed to
	// NewDnsSvcMock before returning. Zero or negative values mean no delay.
	Duration time.Duration
}

// DnsSvcMockImpl is a mock implementation of DnsSvc for testing. It returns
// pre-configured results for each hostname and can simulate lookup delays
// through a TimeSvc. Create instances with NewDnsSvcMock. The zero value
// is not usable. All methods are safe for concurrent use, provided the injected
// TimeSvc is also safe for concurrent calls. However, the returned
// IP slices share the backing array with the registered data (see
// DnsMockResult.IPs), so concurrent mutation of those slices is not safe.
type DnsSvcMockImpl struct {
	timeSvc TimeSvc
	mu      sync.RWMutex
	results map[string]DnsMockResult
}

// NewDnsSvcMock creates a new DnsSvcMockImpl. The timeSvc parameter is used
// to simulate delays for results that have a positive Duration. Pass nil if
// no delay simulation is needed. If timeSvc is nil and a result with a
// positive Duration is looked up, the call panics.
func NewDnsSvcMock(timeSvc TimeSvc) *DnsSvcMockImpl {
	return &DnsSvcMockImpl{
		timeSvc: timeSvc,
		results: make(map[string]DnsMockResult),
	}
}

// LookupIP returns the pre-configured result for host. If no result has been
// registered for host, it returns nil and a *net.DNSError with IsNotFound set
// to true.
//
// The returned IP slice is the same slice that was passed to
// SetLookupIpResult or SetLookupIpResultWithDuration. No copy is made.
func (svc *DnsSvcMockImpl) LookupIP(host string) ([]net.IP, error) {
	// Snapshot the result and release the lock before the optional sleep.
	// Holding the read lock during Sleep would block setters and clears.
	svc.mu.RLock()
	result, exists := svc.results[host]
	svc.mu.RUnlock()

	if !exists {
		return nil, &net.DNSError{
			Err:        "no such host",
			Name:       host,
			Server:     "mock",
			IsNotFound: true,
		}
	}

	if result.Duration > 0 {
		svc.timeSvc.Sleep(result.Duration)
	}

	return result.IPs, result.Err
}

// SetLookupIpResult registers a mock result for host with no delay. The ips
// slice is stored without copying.
func (svc *DnsSvcMockImpl) SetLookupIpResult(host string, ips []net.IP, err error) {
	svc.SetLookupIpResultWithDuration(host, ips, err, 0)
}

// SetLookupIpResultWithDuration registers a mock result for host with an
// optional delay. The ips slice is stored without copying. A positive
// duration causes LookupIP to call timeSvc.Sleep before returning, so
// timeSvc must not be nil.
func (svc *DnsSvcMockImpl) SetLookupIpResultWithDuration(host string, ips []net.IP, err error, duration time.Duration) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.results[host] = DnsMockResult{
		IPs:      ips,
		Err:      err,
		Duration: duration,
	}
}

// ClearLookupIpResult removes the registered result for host. Subsequent
// lookups for this host return the default not-found error.
func (svc *DnsSvcMockImpl) ClearLookupIpResult(host string) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	delete(svc.results, host)
}

// ClearAllLookupIpResults removes all registered results. Subsequent lookups
// for any host return the default not-found error.
func (svc *DnsSvcMockImpl) ClearAllLookupIpResults() {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.results = make(map[string]DnsMockResult)
}
