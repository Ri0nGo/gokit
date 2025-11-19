package maps

import (
	"strconv"
	"sync"
	"testing"
)

// ---------------- ShardMap Benchmark ----------------

// benchmark helper
func benchShardMapSet(b *testing.B, goroutines int) {
	m := NewShardMap[string, int](16)

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i)
			m.Set(key, i)
			i++
		}
	})
}

func benchShardMapGet(b *testing.B, goroutines int) {
	m := NewShardMap[string, int](16)
	for i := 0; i < 1_000_000; i++ {
		m.Set("k"+strconv.Itoa(i), i)
	}

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i&0xFFFF)
			m.Get(key)
			i++
		}
	})
}

func benchShardMapMix(b *testing.B, goroutines int) {
	m := NewShardMap[string, int](16)

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i&0xFFFF)
			if i&1 == 0 {
				m.Set(key, i)
			} else {
				m.Get(key)
			}
			i++
		}
	})
}

// ---------------- sync.Map Benchmark ----------------

func benchSyncMapSet(b *testing.B, goroutines int) {
	var m sync.Map

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i)
			m.Store(key, i)
			i++
		}
	})
}

func benchSyncMapGet(b *testing.B, goroutines int) {
	var m sync.Map
	for i := 0; i < 1_000_000; i++ {
		m.Store("k"+strconv.Itoa(i), i)
	}

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i&0xFFFF)
			m.Load(key)
			i++
		}
	})
}

func benchSyncMapMix(b *testing.B, goroutines int) {
	var m sync.Map

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i&0xFFFF)
			if i&1 == 0 {
				m.Store(key, i)
			} else {
				m.Load(key)
			}
			i++
		}
	})
}

// ---------------- map + RWMutex Benchmark ----------------

func benchMutexMapSet(b *testing.B, goroutines int) {
	m := NewConcurrentMap[string, int]()

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i)
			m.Set(key, i)
			i++
		}
	})
}

func benchMutexMapGet(b *testing.B, goroutines int) {
	m := NewConcurrentMap[string, int]()
	for i := 0; i < 1_000_000; i++ {
		m.Set("k"+strconv.Itoa(i), i)
	}

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i&0xFFFF)
			m.Get(key)
			i++
		}
	})
}

func benchMutexMapMix(b *testing.B, goroutines int) {
	m := NewConcurrentMap[string, int]()

	b.SetParallelism(goroutines)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "k" + strconv.Itoa(i&0xFFFF)
			if i&1 == 0 {
				m.Set(key, i)
			} else {
				m.Get(key)
			}
			i++
		}
	})
}

// ---------------- Public Bench Cases ----------------

func BenchmarkShardMapSet(b *testing.B) { benchShardMapSet(b, 8) }
func BenchmarkShardMapGet(b *testing.B) { benchShardMapGet(b, 8) }
func BenchmarkShardMapMix(b *testing.B) { benchShardMapMix(b, 8) }

func BenchmarkSyncMapSet(b *testing.B) { benchSyncMapSet(b, 8) }
func BenchmarkSyncMapGet(b *testing.B) { benchSyncMapGet(b, 8) }
func BenchmarkSyncMapMix(b *testing.B) { benchSyncMapMix(b, 8) }

func BenchmarkMutexMapSet(b *testing.B) { benchMutexMapSet(b, 8) }
func BenchmarkMutexMapGet(b *testing.B) { benchMutexMapGet(b, 8) }
func BenchmarkMutexMapMix(b *testing.B) { benchMutexMapMix(b, 8) }

// 指定文件执行测试
// go test -bench=. -benchmem -run=^$ shardmap.go concurrentmap.go shardmap_test.go
