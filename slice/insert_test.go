package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	testCase := []struct {
		name      string
		slice     []int
		val       int
		index     int
		wantSlice []int
	}{
		{
			name:      "Index 1",
			slice:     []int{10, 20, 30},
			val:       15,
			index:     1,
			wantSlice: []int{10, 15, 20, 30},
		},
		{
			name:      "Index -1",
			slice:     []int{10, 20, 30},
			val:       0,
			index:     -1,
			wantSlice: nil,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := Insert[int](tc.slice, tc.val, tc.index)
			assert.Equal(t, result, tc.wantSlice)
		})
	}
}

func BenchmarkInsertLow(b *testing.B) {
	s := make([]int, 0)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	for i := 0; i < b.N; i++ {
		InsertLow[int](s, 100, 5)
	}
}

func BenchmarkInsert(b *testing.B) {
	s := make([]int, 0)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	for i := 0; i < b.N; i++ {
		Insert[int](s, 100, 5)
	}
}
