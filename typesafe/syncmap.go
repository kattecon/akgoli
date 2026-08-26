// Package typesafe provides generic wrappers around standard-library
// concurrent data structures.
package typesafe

import "sync"

// SyncMap is a generic wrapper around sync.Map that provides compile-time
// type safety. It exposes a subset of sync.Map operations (Load, Store,
// Delete, LoadOrStore, LoadAndDelete, Range) with typed parameters and
// return values, plus LoadOrCompute for lazy initialization on miss.
// Clear, CompareAndDelete, CompareAndSwap, and Swap are not wrapped.
//
// The zero value is ready for use. A SyncMap must not be copied after first
// use. Like sync.Map, it is optimized for workloads where each entry is
// written once and read many times, or where each goroutine accesses a
// disjoint set of keys. For workloads with frequent writes to shared keys,
// a mutex-guarded map may perform better.
//
// All methods are safe for concurrent use by multiple goroutines. In the
// terminology of the Go memory model, a write operation synchronizes before
// any read operation that observes its effect. Load, LoadAndDelete, and
// LoadOrStore (when loaded is true) are read operations. Store, Delete,
// LoadAndDelete, and LoadOrStore (when loaded is false) are write operations.
// LoadOrCompute follows the same rules as its underlying LoadOrStore call.
// See sync.Map for the full memory model.
type SyncMap[K comparable, V any] struct {
	inner sync.Map
	// noValue holds the zero value of V, returned when a key is not found.
	noValue V
}

// Load returns the value stored under key and true. If the key is absent, it
// returns the zero value of V and false.
func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.inner.Load(key)
	if !ok {
		return m.noValue, ok
	}
	return v.(V), ok
}

// Store sets the value for key.
func (m *SyncMap[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

// Delete removes the value for key.
func (m *SyncMap[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

// LoadOrStore returns the existing value for key if present. Otherwise it
// stores value and returns it. The loaded result is true if the value was
// already present, false if it was stored.
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.inner.LoadOrStore(key, value)
	return a.(V), loaded
}

// LoadOrCompute returns the existing value for key if present. On a miss, it
// calls f to compute a new value, stores it with LoadOrStore, and returns it.
//
// Under concurrent access, multiple goroutines may miss simultaneously and
// each call f. Only one computed value is kept in the map. The others are
// discarded. Callers must ensure that f is safe to run more than once and
// does not leak resources on the losing calls. The loaded result is true when
// an existing value was returned, false when this caller's computed value was
// stored.
func (m *SyncMap[K, V]) LoadOrCompute(key K, f func() V) (actual V, loaded bool) {
	// Try Load first to avoid running f() when the key already exists.
	// On a miss this costs two map lookups instead of one.

	a, loaded := m.Load(key)
	if loaded {
		return a, loaded
	}

	return m.LoadOrStore(key, f())
}

// LoadAndDelete removes the value for key and returns the previous value if
// any. The loaded result is true if the key was present, false otherwise. When
// false, the returned value is the zero value of V.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	a, loaded := m.inner.LoadAndDelete(key)
	if !loaded {
		return m.noValue, loaded
	}
	return a.(V), loaded
}

// Range calls f sequentially for each key-value pair in the map. If f returns
// false, Range stops the iteration. Range does not block other map methods,
// and f may call any method on the same map without deadlocking. Range does
// not provide a consistent snapshot of the map. A key may be visited once,
// skipped, or reflect a concurrent update. Range takes O(N) time even if f
// returns false after a small number of calls.
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
