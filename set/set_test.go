package set

import (
	"reflect"
	"sort"
	"testing"
)

func TestSetBasicOperations(t *testing.T) {
	s := NewSet[int]()

	// Add 元素
	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}

	// Contains
	if !s.Contains(2) {
		t.Errorf("expected set to contain 2")
	}
	if s.Contains(4) {
		t.Errorf("expected set not to contain 4")
	}

	// Delete
	s.Delete(2)
	if s.Len() != 2 || s.Contains(2) {
		t.Errorf("delete failed, len=%d, contains 2=%v", s.Len(), s.Contains(2))
	}

	// Clear
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("clear failed, len=%d", s.Len())
	}
}

func TestSetItems(t *testing.T) {
	s := NewSet[string]()
	s.Add("a", "b", "c")

	items := s.Items()
	sort.Strings(items) // 保证顺序可比
	expected := []string{"a", "b", "c"}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("expected items %v, got %v", expected, items)
	}
}

func TestSetUnion(t *testing.T) {
	a := NewSet[int]()
	b := NewSet[int]()

	a.Add(1, 2, 3)
	b.Add(3, 4, 5)

	u := a.Union(b)
	expected := NewSet[int]()
	expected.Add(1, 2, 3, 4, 5)

	if !compareSets(u, expected) {
		t.Errorf("union failed, got %v", u.Items())
	}
}

func TestSetIntersect(t *testing.T) {
	a := NewSet[int]()
	b := NewSet[int]()

	a.Add(1, 2, 3)
	b.Add(2, 3, 4)

	inter := a.Intersect(b)
	expected := NewSet[int]()
	expected.Add(2, 3)

	if !compareSets(inter, expected) {
		t.Errorf("intersect failed, got %v", inter.Items())
	}
}

func TestSetDifference(t *testing.T) {
	a := NewSet[int]()
	b := NewSet[int]()

	a.Add(1, 2, 3, 4)
	b.Add(3, 4, 5)

	diff := a.Difference(b)
	expected := NewSet[int]()
	expected.Add(1, 2)

	if !compareSets(diff, expected) {
		t.Errorf("difference failed, got %v", diff.Items())
	}
}

// 辅助函数：比较两个 set 是否相等
func compareSets[T comparable](a, b *Set[T]) bool {
	if a.Len() != b.Len() {
		return false
	}
	for _, v := range a.Items() {
		if !b.Contains(v) {
			return false
		}
	}
	return true
}
