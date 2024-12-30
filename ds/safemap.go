package ds

import (
	"fmt"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	container map[K]V
	mu        sync.RWMutex
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		container: make(map[K]V),
	}
}

func (m *SafeMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.container[k] = v
}

func (m *SafeMap[K, V]) Sets(em map[K]V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range em {
		m.container[k] = v
	}
}

func (m *SafeMap[K, V]) Del(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.container, k)
}

func (m *SafeMap[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.container[k]
	return v, ok
}

func (m *SafeMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, 0, len(m.container))
	for k := range m.container {
		keys = append(keys, k)
	}

	return keys
}

func (m *SafeMap[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return MapValues(m.container)
}

func MapValues[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for k := range m {
		result = append(result, m[k])
	}
	return result
}

func (m *SafeMap[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.container)
}

func (m *SafeMap[K, V]) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return fmt.Sprintf("%#v", m.container)
}
