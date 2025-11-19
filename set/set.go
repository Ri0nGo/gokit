package set

// 通过map实现的set数据结构
// 不是并发安全的！！！

type Set[T comparable] struct {
	container map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		container: make(map[T]struct{}),
	}
}

// Add 添加元素
func (s *Set[T]) Add(elems ...T) {
	for _, elem := range elems {
		s.container[elem] = struct{}{}
	}
}

// Delete 删除元素
func (s *Set[T]) Delete(elem T) {
	delete(s.container, elem)
}

// Clear 清空set
func (s *Set[T]) Clear() {
	s.container = make(map[T]struct{})
}

// Len 统计set长度
func (s *Set[T]) Len() int {
	return len(s.container)
}

// Contains 是否包含元素
func (s *Set[T]) Contains(elem T) bool {
	_, ok := s.container[elem]
	return ok
}

// Items 返回set中的所有元素
func (s *Set[T]) Items() []T {
	keys := make([]T, 0, len(s.container))
	for key := range s.container {
		keys = append(keys, key)
	}
	return keys
}

// Union 并集
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, v := range s.Items() {
		result.Add(v)
	}
	for _, v := range other.Items() {
		result.Add(v)
	}
	return result
}

// Intersect 交集
func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, v := range s.Items() {
		if other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference 差集
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, v := range s.Items() {
		if !other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}
