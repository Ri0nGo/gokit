package maps

import "sync"

// 基于sync.RWMutex 实现的一个并发安全的map

type ConcurrentMap[K comparable, V any] struct {
	container map[K]V
	mux       sync.RWMutex
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		container: make(map[K]V),
	}
}

func (m *ConcurrentMap[K, V]) Set(key K, val V) {
	m.mux.Lock()
	m.container[key] = val
	m.mux.Unlock()
}

func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	v, ok := m.container[key]
	return v, ok
}

func (m *ConcurrentMap[K, V]) Delete(key K) {
	m.mux.Lock()
	delete(m.container, key)
	m.mux.Unlock()
}

func (m *ConcurrentMap[K, V]) Len() int {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return len(m.container)
}

func (m *ConcurrentMap[K, V]) Keys() []K {
	m.mux.RLock()
	keys := make([]K, 0, len(m.container))
	for key := range m.container {
		keys = append(keys, key)
	}
	m.mux.RUnlock()
	return keys
}

func (m *ConcurrentMap[K, V]) Values() []V {
	m.mux.RLock()
	defer m.mux.RUnlock()

	values := make([]V, 0, len(m.container))
	for _, v := range m.container {
		values = append(values, v)
	}
	return values
}
