package absos

import (
	"runtime"
	"time"
)

type TimeSvcMockImpl struct {
	Time time.Time

	StopSleepingRequest chan any
	DoneSleeping        chan time.Duration
}

func NewTimeSvcMock() *TimeSvcMockImpl {
	return &TimeSvcMockImpl{
		StopSleepingRequest: make(chan any),
		DoneSleeping:        make(chan time.Duration),
	}
}

func (svc *TimeSvcMockImpl) Now() time.Time {
	return svc.Time
}

func (svc *TimeSvcMockImpl) Add(d time.Duration) {
	svc.Time = svc.Time.Add(d)
}

func (svc *TimeSvcMockImpl) Sleep(d time.Duration) {
	<-svc.StopSleepingRequest
	svc.DoneSleeping <- d
	runtime.Gosched()
}

func (svc *TimeSvcMockImpl) GetSleptDuration() time.Duration {
	svc.StopSleepingRequest <- nil
	return <-svc.DoneSleeping
}
