package typesafe

import "sync"

type SyncMap[K comparable, V any] struct {
	inner   sync.Map
	noValue V
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.inner.Load(key)
	if !ok {
		return m.noValue, ok
	}
	return v.(V), ok
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.inner.LoadOrStore(key, value)
	return a.(V), loaded
}

func (m *SyncMap[K, V]) LoadOrCompute(key K, f func() V) (actual V, loaded bool) {
	// Assuming misses are rare.

	a, loaded := m.Load(key)
	if loaded {
		return a, loaded
	}

	return m.LoadOrStore(key, f())
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	a, loaded := m.inner.LoadAndDelete(key)
	if !loaded {
		return m.noValue, loaded
	}
	return a.(V), loaded
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}
