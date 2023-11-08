package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
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
			result, _ := Add[int](tc.slice, tc.val, tc.index)
			assert.Equal(t, result, tc.wantSlice)
		})
	}
}

func BenchmarkAddLow(b *testing.B) {
	s := make([]int, 0)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	for i := 0; i < b.N; i++ {
		AddLow[int](s, 100, 5)
	}
}

func BenchmarkAdd(b *testing.B) {
	s := make([]int, 0)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	for i := 0; i < b.N; i++ {
		Add[int](s, 100, 5)
	}
}

// -- 测试结果 -- //
/*
go test -bench=. -benchmem

BenchmarkAdd-16         32219262                40.33 ns/op           96 B/op          1 allocs/op
BenchmarkAdd2-16        237656610                4.288 ns/op           0 B/op          0 allocs/op
*/

// 参考链接： https://juejin.cn/post/6970615934255906830
