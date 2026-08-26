package absos

import "time"

// TimeSvc abstracts the system clock. The production implementation delegates
// to time.Now and time.Sleep. Use NewTimeSvcMock for a controlled test
// replacement that lets tests advance time without real delays.
type TimeSvc interface {
	// Now returns the current time. The production implementation calls
	// time.Now.
	Now() time.Time

	// Sleep pauses the current goroutine for at least duration d. The
	// production implementation calls time.Sleep, which returns immediately
	// for non-positive durations.
	Sleep(d time.Duration)
}

type timeSvcImpl struct{}

var timeSvcImplInstance = timeSvcImpl{}

// NewTimeSvc returns the production TimeSvc that delegates to time.Now and
// time.Sleep. The returned value is stateless and shared across callers.
func NewTimeSvc() TimeSvc {
	return timeSvcImplInstance
}

func (timeSvcImpl) Now() time.Time {
	return time.Now()
}

func (timeSvcImpl) Sleep(d time.Duration) {
	time.Sleep(d)
}
