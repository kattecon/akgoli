package absos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeSvc(t *testing.T) {
	svc := NewTimeSvc()

	assert.Equal(t, svc, NewTimeSvc())

	t1 := svc.Now()
	time.Sleep(1 * time.Millisecond)
	t2 := svc.Now()

	assert.True(t, t2.After(t1))
}

func TestTimeSvcSleep(t *testing.T) {
	svc := NewTimeSvc()

	t1 := svc.Now()
	svc.Sleep(10 * time.Millisecond)
	t2 := svc.Now()

	assert.Greater(t, t2.Sub(t1), 9*time.Microsecond)
}
