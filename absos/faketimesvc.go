package absos

import (
	"runtime"
	"time"
)

type FakeTimeSvcImpl struct {
	Time time.Time

	StopSleepingRequest chan any
	DoneSleeping        chan time.Duration
}

func NewFakeTimeSvc() *FakeTimeSvcImpl {
	return &FakeTimeSvcImpl{
		StopSleepingRequest: make(chan any),
		DoneSleeping:        make(chan time.Duration),
	}
}

func (svc *FakeTimeSvcImpl) Now() time.Time {
	return svc.Time
}

func (svc *FakeTimeSvcImpl) Add(d time.Duration) {
	svc.Time = svc.Time.Add(d)
}

func (svc *FakeTimeSvcImpl) Sleep(d time.Duration) {
	<-svc.StopSleepingRequest
	svc.DoneSleeping <- d
	runtime.Gosched()
}

func (svc *FakeTimeSvcImpl) GetSleptDuration() time.Duration {
	svc.StopSleepingRequest <- nil
	return <-svc.DoneSleeping
}
