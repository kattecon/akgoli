package absos

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDnsSvcMock(t *testing.T) {
	timeSvc := NewTimeSvcMock()
	dnsSvc := NewDnsSvcMock(timeSvc)

	// Test unknown host returns error.
	_, err := dnsSvc.LookupIP("unknown-host.com")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such host")

	// Test setting and getting results.
	expectedIPs := []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("192.168.1.2")}
	dnsSvc.SetLookupIpResult("example.com", expectedIPs, nil)

	ips, err := dnsSvc.LookupIP("example.com")
	assert.Nil(t, err)
	assert.Equal(t, expectedIPs, ips)

	// Test setting error result.
	testErr := &net.DNSError{Err: "test error", Name: "error.com"}
	dnsSvc.SetLookupIpResult("error.com", nil, testErr)

	ips, err = dnsSvc.LookupIP("error.com")
	assert.Nil(t, ips)
	assert.Equal(t, testErr, err)
}

func TestDnsSvcMockWithDuration(t *testing.T) {
	timeSvc := NewTimeSvcMock()
	dnsSvc := NewDnsSvcMock(timeSvc)

	expectedIPs := []net.IP{net.ParseIP("10.0.0.1")}
	duration := 100 * time.Millisecond

	// Set result with duration.
	dnsSvc.SetLookupIpResultWithDuration("slow.com", expectedIPs, nil, duration)

	// Start lookup in goroutine.
	var resultIPs []net.IP
	var resultErr error
	done := make(chan bool)

	go func() {
		resultIPs, resultErr = dnsSvc.LookupIP("slow.com")
		done <- true
	}()

	// Wait for the DNS lookup to register as a sleeper.
	timeSvc.WaitForSleepers(1)

	// Verify time hasn't advanced yet.
	initialTime := timeSvc.Now()

	// Advance to the sleep event.
	advanced := timeSvc.AdvanceToNextSleepEvent()
	assert.Equal(t, duration, advanced)

	// Wait for completion.
	<-done

	// Verify results.
	assert.Nil(t, resultErr)
	assert.Equal(t, expectedIPs, resultIPs)
	assert.Equal(t, initialTime.Add(duration), timeSvc.Now())
}

func TestDnsSvcMockClearResults(t *testing.T) {
	timeSvc := NewTimeSvcMock()
	dnsSvc := NewDnsSvcMock(timeSvc)

	// Set some results.
	dnsSvc.SetLookupIpResult("host1.com", []net.IP{net.ParseIP("1.1.1.1")}, nil)
	dnsSvc.SetLookupIpResult("host2.com", []net.IP{net.ParseIP("2.2.2.2")}, nil)

	// Verify they work.
	ips, err := dnsSvc.LookupIP("host1.com")
	assert.Nil(t, err)
	assert.Equal(t, []net.IP{net.ParseIP("1.1.1.1")}, ips)

	// Clear specific result.
	dnsSvc.ClearLookupIpResult("host1.com")

	// host1.com should now fail.
	_, err = dnsSvc.LookupIP("host1.com")
	assert.NotNil(t, err)

	// host2.com should still work.
	ips, err = dnsSvc.LookupIP("host2.com")
	assert.Nil(t, err)
	assert.Equal(t, []net.IP{net.ParseIP("2.2.2.2")}, ips)

	// Clear all results.
	dnsSvc.ClearAllLookupIpResults()

	// Both should now fail.
	_, err = dnsSvc.LookupIP("host2.com")
	assert.NotNil(t, err)
}

func TestDnsSvcMockZeroDuration(t *testing.T) {
	timeSvc := NewTimeSvcMock()
	dnsSvc := NewDnsSvcMock(timeSvc)

	expectedIPs := []net.IP{net.ParseIP("127.0.0.1")}

	// Set result with zero duration.
	dnsSvc.SetLookupIpResultWithDuration("fast.com", expectedIPs, nil, 0)

	initialTime := timeSvc.Now()

	// Should return immediately without sleeping.
	ips, err := dnsSvc.LookupIP("fast.com")
	assert.Nil(t, err)
	assert.Equal(t, expectedIPs, ips)

	// Time should not have advanced.
	assert.Equal(t, initialTime, timeSvc.Now())
}
