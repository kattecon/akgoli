package absos

import (
	"runtime"
	"sync"
	"time"
)

// sleepRequest represents a single Sleep() call with its due time and synchronization channels.
type sleepRequest struct {
	dueTime     time.Time
	releaseChan chan any // Used to signal the sleeper to wake up.
	doneChan    chan any // Used by sleeper to signal completion.
}

// TimeSvcMockImpl provides a mock implementation of TimeSvc for testing purposes.
// It allows controlled time advancement and synchronization with sleeping goroutines.
//
// The mock works by intercepting Sleep() calls and coordinating them with test code
// via AdvanceToNextSleepEvent() or Add(). This enables deterministic testing of time-dependent
// code without relying on real time delays.
//
// Thread Safety: All methods are thread-safe and can be called from multiple goroutines.
// However, the exported Time field is not protected by the internal mutex. Use Now
// to read the current time during concurrent work. Direct reads or writes of Time
// while another goroutine calls Sleep, Add, or AdvanceToNextSleepEvent cause a data race.
//
// Usage Pattern:
//  1. Start a goroutine that calls Sleep.
//  2. Call WaitForSleepers to block until the goroutine has registered.
//  3. Call Add or AdvanceToNextSleepEvent to advance the clock and release the sleeper.
//
// Skipping step 2 is a race: the sleeper may register after Add has already run,
// leaving the sleeper blocked forever.
type TimeSvcMockImpl struct {
	// mu protects all fields below from concurrent access.
	mu sync.Mutex

	// Time represents the current mock time. Read it directly only when no other
	// goroutine is calling Sleep, Add, or AdvanceToNextSleepEvent. Use Now for
	// safe concurrent reads.
	Time time.Time

	// sleepers contains all currently sleeping goroutines waiting to be awakened.
	sleepers []*sleepRequest
}

// NewTimeSvcMock creates a new TimeSvcMockImpl instance.
//
// The returned mock starts with zero time and no sleepers.
func NewTimeSvcMock() *TimeSvcMockImpl {
	return &TimeSvcMockImpl{
		sleepers: make([]*sleepRequest, 0),
	}
}

// SleeperCount returns the current number of registered sleepers.
//
// This method is useful for testing to ensure all expected sleepers
// have registered before proceeding with time advancement.
// Thread-safe for concurrent access.
func (svc *TimeSvcMockImpl) SleeperCount() int {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	return len(svc.sleepers)
}

// WaitForSleepers waits until the specified number of sleepers are registered.
//
// This method provides deterministic synchronization for tests, ensuring
// all expected Sleep() calls have registered before advancing time.
// There is no timeout. If the requested count is never reached, the call
// blocks indefinitely. The method polls with a real 1ms time.Sleep between
// checks.
// Thread-safe for concurrent access.
func (svc *TimeSvcMockImpl) WaitForSleepers(count int) {
	for {
		svc.mu.Lock()
		currentCount := len(svc.sleepers)
		svc.mu.Unlock()

		if currentCount >= count {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
}

// Now returns the current mock time.
//
// This method is part of the TimeSvc interface.
// This time does not advance automatically - it only changes when Add() or
// AdvanceToNextSleepEvent() is called.
// Thread-safe for concurrent access.
func (svc *TimeSvcMockImpl) Now() time.Time {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	return svc.Time
}

// Add advances the mock time by the specified duration.
//
// This method automatically releases all sleepers whose due time has been reached
// or passed after the time advancement.
// Thread-safe for concurrent access.
func (svc *TimeSvcMockImpl) Add(d time.Duration) {
	// Lock scoped to the anonymous function so it releases before
	// releaseReadySleepers, which re-acquires it.
	func() {
		svc.mu.Lock()
		defer svc.mu.Unlock()
		svc.Time = svc.Time.Add(d)
	}()

	svc.releaseReadySleepers()
}

// releaseReadySleepers finds and releases all sleepers whose due time has been reached.
// This method must be called without holding svc.mu as it will acquire and release the lock itself.
func (svc *TimeSvcMockImpl) releaseReadySleepers() {
	var readySleepers []*sleepRequest

	func() {
		svc.mu.Lock()
		defer svc.mu.Unlock()

		// Find sleepers whose due time has been reached.
		var remaining []*sleepRequest
		for _, sleeper := range svc.sleepers {
			if svc.Time.Before(sleeper.dueTime) {
				// This sleeper's due time has not been reached yet.
				remaining = append(remaining, sleeper)
			} else {
				// This sleeper's due time has been reached.
				readySleepers = append(readySleepers, sleeper)
			}
		}
		svc.sleepers = remaining
	}()

	// Release sleepers without holding the lock.
	for _, sleeper := range readySleepers {
		sleeper.releaseChan <- nil

		// Wait for sleeper to actually wake up and complete.
		<-sleeper.doneChan
	}
}

// Sleep simulates sleeping for the specified duration.
//
// The method calculates the due time (current time + duration) and blocks
// until that time is reached via Add() or AdvanceToNextSleepEvent().
// Thread-safe: Multiple goroutines can call Sleep() concurrently.
func (svc *TimeSvcMockImpl) Sleep(d time.Duration) {
	// Buffer of 1 lets the sleeper proceed past the doneChan send without
	// waiting for releaseReadySleepers to receive. The subsequent Gosched
	// yield gives the releaser a chance to observe the send.
	releaseChan := make(chan any, 1)
	doneChan := make(chan any, 1)

	// Lock scoped to the anonymous function so it releases before blocking on
	// releaseChan. Holding the lock during the channel receive would deadlock:
	// Add could never acquire the lock to advance the clock and release this
	// sleeper.
	func() {
		svc.mu.Lock()
		defer svc.mu.Unlock()

		dueTime := svc.Time.Add(d)
		request := &sleepRequest{
			dueTime:     dueTime,
			releaseChan: releaseChan,
			doneChan:    doneChan,
		}
		svc.sleepers = append(svc.sleepers, request)
	}()

	// Block until released by Add() or AdvanceToNextSleepEvent().
	<-releaseChan

	// Signal completion and yield. The yield gives the releasing goroutine
	// (releaseReadySleepers) a chance to proceed before this goroutine
	// continues, preventing ordering-sensitive tests from seeing stale
	// sleeper counts.
	doneChan <- nil
	runtime.Gosched()
}

// AdvanceToNextSleepEvent advances mock time to the next sleep event and releases those sleepers.
//
// This method finds the minimum due time among all current sleepers, advances
// the mock time to that point, and releases all sleepers whose due time has
// been reached.
//
// Returns 0 if no sleepers are present, otherwise returns the duration advanced.
// Thread-safe for concurrent access.
func (svc *TimeSvcMockImpl) AdvanceToNextSleepEvent() time.Duration {
	var currentTime time.Time
	var minDueTime time.Time
	var found bool

	// Lock scoped to the anonymous function so it releases before
	// releaseReadySleepers, which re-acquires it.
	func() {
		svc.mu.Lock()
		defer svc.mu.Unlock()

		currentTime = svc.Time

		// Find the minimum due time among all sleepers.
		for _, sleeper := range svc.sleepers {
			if !found || sleeper.dueTime.Before(minDueTime) {
				minDueTime = sleeper.dueTime
				found = true
			}
		}

		if found {
			// Advance time to the minimum due time.
			svc.Time = minDueTime
		}
	}()

	if found {
		svc.releaseReadySleepers()
		return minDueTime.Sub(currentTime)
	}
	return 0
}
