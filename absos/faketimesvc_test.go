package absos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFakeTimeSvc(t *testing.T) {
	svc := NewFakeTimeSvc()

	t1 := svc.Now()
	time.Sleep(1 * time.Millisecond)
	t2 := svc.Now()

	svc.Add(100 * time.Microsecond)

	assert.NotSame(t, t1, t2)
	assert.Equal(t, t1, t2)
	assert.Equal(t, svc.Time, svc.Now())
	assert.Equal(t, t2.Add(100*time.Microsecond), svc.Now())
}

func TestFakeTimeSvcSleep(t *testing.T) {
	svc := NewFakeTimeSvc()

	go func() {
		svc.Sleep(100)
		svc.Sleep(300)
	}()

	assert.Equal(t, time.Duration(100), svc.GetSleptDuration())
	assert.Equal(t, time.Duration(300), svc.GetSleptDuration())
}
