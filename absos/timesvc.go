package absos

import "time"

type TimeSvc interface {
	Now() time.Time
	Sleep(d time.Duration)
}

type TimeSvcImpl struct{}

var timeSvcImpl = TimeSvcImpl{}

func NewTimeSvc() TimeSvc {
	return timeSvcImpl
}

func (TimeSvcImpl) Now() time.Time {
	return time.Now()
}

func (TimeSvcImpl) Sleep(d time.Duration) {
	time.Sleep(d)
}
