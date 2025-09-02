package absos

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeSvcMock(t *testing.T) {
	svc := NewTimeSvcMock()

	t1 := svc.Now()
	time.Sleep(1 * time.Millisecond)
	t2 := svc.Now()

	svc.Add(100 * time.Microsecond)

	assert.NotSame(t, &t1, &t2)
	assert.Equal(t, t1, t2)
	assert.Equal(t, svc.Time, svc.Now())
	assert.Equal(t, t2.Add(100*time.Microsecond), svc.Now())
}

func TestTimeSvcMockSleep(t *testing.T) {
	svc := NewTimeSvcMock()

	go func() {
		svc.Sleep(100 * time.Nanosecond)
		svc.Sleep(300 * time.Nanosecond)
	}()

	svc.WaitForSleepers(1)

	// Advance to next sleep event and verify the duration.
	advanced := svc.AdvanceToNextSleepEvent()
	assert.Equal(t, 100*time.Nanosecond, advanced)

	svc.WaitForSleepers(1)

	advanced = svc.AdvanceToNextSleepEvent()
	assert.Equal(t, 300*time.Nanosecond, advanced)
}

func TestTimeSvcMockMultipleGoroutines(t *testing.T) {
	svc := NewTimeSvcMock()
	var wg sync.WaitGroup
	const numGoroutines = 5

	// Start multiple goroutines with different sleep durations.
	durations := []time.Duration{
		500 * time.Nanosecond,
		200 * time.Nanosecond,
		800 * time.Nanosecond,
		100 * time.Nanosecond,
		400 * time.Nanosecond,
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(d time.Duration) {
			defer wg.Done()
			svc.Sleep(d)
		}(durations[i])
	}

	// Wait for all sleepers to register before proceeding.
	svc.WaitForSleepers(numGoroutines)

	// Advance to each sleep event in order of due time.
	// Expected order: 100, 200, 400, 500, 800 nanoseconds.
	expectedAdvances := []time.Duration{
		100 * time.Nanosecond, // 0 → 100
		100 * time.Nanosecond, // 100 → 200
		200 * time.Nanosecond, // 200 → 400
		100 * time.Nanosecond, // 400 → 500
		300 * time.Nanosecond, // 500 → 800
	}

	for i, expected := range expectedAdvances {
		svc.WaitForSleepers(1)
		advanced := svc.AdvanceToNextSleepEvent()
		assert.Equal(t, expected, advanced, "Advance step %d", i+1)
	}

	// Wait for all goroutines to complete.
	wg.Wait()

	// Verify final time.
	assert.Equal(t, 800*time.Nanosecond, svc.Now().Sub(time.Time{}))
}

func TestTimeSvcMockAddReleasesReadySleepers(t *testing.T) {
	svc := NewTimeSvcMock()
	var completed sync.WaitGroup

	// Start a goroutine that sleeps for 100 nanoseconds.
	completed.Add(1)
	go func() {
		defer completed.Done()
		svc.Sleep(100 * time.Nanosecond)
	}()

	// Wait for the sleeper to register.
	svc.WaitForSleepers(1)

	// Add 150 nanoseconds - this should release the sleeper.
	svc.Add(150 * time.Nanosecond)

	// The goroutine should complete now.
	completed.Wait()

	// Verify the time was advanced.
	assert.Equal(t, 150*time.Nanosecond, svc.Now().Sub(time.Time{}))
}

func TestTimeSvcMockPartialAdvancement(t *testing.T) {
	svc := NewTimeSvcMock()
	var wg sync.WaitGroup

	// Start sleepers with different durations.
	wg.Add(2)

	go func() {
		defer wg.Done()
		svc.Sleep(100 * time.Nanosecond)
	}()

	go func() {
		defer wg.Done()
		svc.Sleep(200 * time.Nanosecond)
	}()

	// Wait for sleepers to register.
	svc.WaitForSleepers(2)

	// Advance by 150 nanoseconds - should only release the first sleeper.
	svc.Add(150 * time.Nanosecond)

	// There still should be 1 sleeper.
	svc.WaitForSleepers(1)

	// Advance to the next sleep event.
	advanced := svc.AdvanceToNextSleepEvent()
	assert.Equal(t, 50*time.Nanosecond, advanced) // 200 - 150

	// Make sure there are no sleepers left.
	assert.Equal(t, 0, svc.SleeperCount())

	wg.Wait()
}
