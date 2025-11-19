package maps

import (
	"sync"
	"sync/atomic"
	"testing"
)

const (
	ops = 10000
)

func BenchmarkHandWriteSafe(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	var wg sync.WaitGroup
	var counter atomic.Int64

	b.ResetTimer()
	for i := 0; i < b.N; i++ { // 支持 benchmark 多次运行
		wg.Add(2)

		// writer
		go func() {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				m.Set(j, j*10)
				if j%7 == 0 {
					m.Delete(j / 2)
				}
			}
		}()

		// reader
		go func() {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				if v, ok := m.Get(j % 1000); ok {
					counter.Add(int64(v))
				}
			}
		}()

		wg.Wait()
	}
}

func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	var wg sync.WaitGroup
	var counter atomic.Int64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)

		go func() {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				m.Store(j, j*10)
				if j%7 == 0 {
					m.Delete(j / 2)
				}
			}
		}()

		go func() {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				if v, ok := m.Load(j % 1000); ok {
					if iv, _ := v.(int); iv > 0 {
						counter.Add(int64(iv))
					}
				}
			}
		}()

		wg.Wait()
	}
}
