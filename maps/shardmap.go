package maps

import (
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
)

var defaultShardCnt uint64 = 16

type shard[K comparable, V any] struct {
	container map[K]V
	mu        sync.RWMutex
}

type ShardMap[K comparable, V any] struct {
	shards    []*shard[K, V]
	hasher    func(K) uint64
	count     uint64
	shardMask uint64
	total     atomic.Uint64
}

// NewShardMap 创建分片数为 >= shardCnt 的最小 2^n（若传入 <=0，则使用默认值）
func NewShardMap[K comparable, V any](shardCnt uint64) *ShardMap[K, V] {
	if shardCnt <= 0 {
		shardCnt = defaultShardCnt
	}
	// 向上取最接近的 2^n
	count := roundUpToPower2(shardCnt)

	sm := &ShardMap[K, V]{
		shards:    make([]*shard[K, V], count),
		hasher:    defaultHasher[K],
		count:     count,
		shardMask: count - 1,
	}
	// 初始化分段后的map
	for i := 0; i < int(count); i++ {
		sm.shards[i] = &shard[K, V]{
			container: make(map[K]V),
		}
	}
	return sm
}

// Set 设置 key, value
func (s *ShardMap[K, V]) Set(key K, val V) {
	shardMap := s.getShard(key)
	shardMap.mu.Lock()

	if _, existed := shardMap.container[key]; !existed {
		s.total.Add(1)
	}
	shardMap.container[key] = val
	shardMap.mu.Unlock()
}

// Get 获取值
func (s *ShardMap[K, V]) Get(key K) (V, bool) {
	shardMap := s.getShard(key)
	shardMap.mu.RLock()
	val, ok := shardMap.container[key]
	shardMap.mu.RUnlock()
	return val, ok
}

// Delete 删除 key
func (s *ShardMap[K, V]) Delete(key K) {
	shardMap := s.getShard(key)
	shardMap.mu.Lock()
	if _, existed := shardMap.container[key]; existed {
		s.total.Add(^uint64(0)) // -1 using atomic.Uint64 (Add with wrap); or use Add(^uint64(0)) equals -1
		// Alternatively: s.total.Add(^uint64(0)) is a trick; clearer:
		// s.total.Add(^uint64(0)) // subtract 1
		// But for readability, use s.total.Add(^uint64(0))
	}
	delete(shardMap.container, key)
	shardMap.mu.Unlock()
}

// Len 返回近似总元素个数（极快，几乎无锁）
func (s *ShardMap[K, V]) Len() uint64 {
	return s.total.Load()
}

// Keys 返回所有 key（会分配内存，但比全局锁快很多）
func (s *ShardMap[K, V]) Keys() []K {
	// 预分配近似容量
	keys := make([]K, 0, s.Len())

	for _, sm := range s.shards {
		sm.mu.RLock()
		for k := range sm.container {
			keys = append(keys, k)
		}
		sm.mu.RUnlock()
	}
	return keys
}

// Values 返回所有 value
func (s *ShardMap[K, V]) Values() []V {
	values := make([]V, 0, s.Len())
	for _, sm := range s.shards {
		sm.mu.RLock()
		for _, v := range sm.container {
			values = append(values, v)
		}
		sm.mu.RUnlock()
	}
	return values
}

// Range 遍历所有键值, f返回true则停止。通过 snapshot（复制）实现，
// 回调中可以安全调用 Set/Delete（但不保证写入能被当前 Range 看到）
func (s *ShardMap[K, V]) Range(f func(key K, value V) (stop bool)) {
	type pair struct {
		k K
		v V
	}
	items := make([]pair, 0, s.Len())

	for _, sh := range s.shards {
		sh.mu.RLock()
		for k, v := range sh.container {
			items = append(items, pair{k, v})
		}
		sh.mu.RUnlock()
	}

	for _, it := range items {
		if f(it.k, it.v) {
			return
		}
	}
}

// ---------------- 辅助函数 ---------------- //

// 获取 key 对应的 shard
func (s *ShardMap[K, V]) getShard(key K) *shard[K, V] {
	idx := s.hasher(key) & s.shardMask
	return s.shards[idx]
}

// defaultHasher 指向 fastHash，便于以后替换
func defaultHasher[K comparable](k K) uint64 {
	return fastHash(k)
}

// fastHash 对常见类型做特化，fallback 到 fmt+fnv
func fastHash[K comparable](k K) uint64 {
	switch x := any(k).(type) {
	case string:
		// FNV-1a inline (no allocation)
		var h uint64 = 1469598103934665603
		for i := 0; i < len(x); i++ {
			h ^= uint64(x[i])
			h *= 1099511628211
		}
		return h
	case int:
		u := uint64(x)
		u = (u ^ (u >> 30)) * 0xbf58476d1ce4e5b9
		u = (u ^ (u >> 27)) * 0x94d049bb133111eb
		return u ^ (u >> 31)
	case uint:
		u := uint64(x)
		u = (u ^ (u >> 30)) * 0xbf58476d1ce4e5b9
		u = (u ^ (u >> 27)) * 0x94d049bb133111eb
		return u ^ (u >> 31)
	case int64:
		u := uint64(x)
		u = (u ^ (u >> 30)) * 0xbf58476d1ce4e5b9
		u = (u ^ (u >> 27)) * 0x94d049bb133111eb
		return u ^ (u >> 31)
	case uint64:
		u := uint64(x)
		u = (u ^ (u >> 30)) * 0xbf58476d1ce4e5b9
		u = (u ^ (u >> 27)) * 0x94d049bb133111eb
		return u ^ (u >> 31)
	// 根据需要可以继续添加 int32/uint32/[]byte 等特化
	default:
		s := fmt.Sprintf("%v", k)
		h := fnv.New64a()
		_, _ = h.Write([]byte(s))
		return h.Sum64()
	}
}

// roundUpToPower2 向上取最近的 2^n（传入非 0）
func roundUpToPower2(v uint64) uint64 {
	if v == 0 {
		return 1
	}
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}
