package absos

import "time"

type TimeSvc interface {
	Now() time.Time
	Sleep(d time.Duration)
}

type timeSvcImpl struct{}

var timeSvcImplInstance = timeSvcImpl{}

func NewTimeSvc() TimeSvc {
	return timeSvcImplInstance
}

func (timeSvcImpl) Now() time.Time {
	return time.Now()
}

func (timeSvcImpl) Sleep(d time.Duration) {
	time.Sleep(d)
}
