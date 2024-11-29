package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func generateIntSliceByLen(len int) []int {
	var result []int
	if len < 0 {
		return result
	}
	for i := 0; i < len; i++ {
		result = append(result, i)
	}
	return result
}

// TestSplitChunk 单元测试
func TestSplitChunk(t *testing.T) {
	testCase := []struct {
		name         string
		slice        []int
		splitSize    int
		wantLen      int
		lastSliceLen int
		wantErr      error
	}{
		{
			name:         "normal split slice",
			slice:        generateIntSliceByLen(28),
			splitSize:    6,
			wantLen:      5,
			lastSliceLen: 4,
		},
		{
			name:         "full split slice",
			slice:        generateIntSliceByLen(30),
			splitSize:    30,
			wantLen:      1,
			lastSliceLen: 30,
		},
		{
			name:         "one split slice",
			slice:        generateIntSliceByLen(101),
			splitSize:    1,
			wantLen:      101,
			lastSliceLen: 1,
		},
		{
			name:         "splitSize is zero",
			slice:        generateIntSliceByLen(4),
			splitSize:    0,
			wantLen:      0,
			lastSliceLen: 0,
			wantErr:      ErrSplitSizeOutOfSlice,
		},
		{
			name:         "splitSize more than length",
			slice:        generateIntSliceByLen(124),
			splitSize:    125,
			wantLen:      0,
			lastSliceLen: 0,
			wantErr:      ErrSplitSizeOutOfSlice,
		},
		{
			name:         "slice is zero, splitSize is zero",
			slice:        generateIntSliceByLen(0),
			splitSize:    0,
			wantLen:      0,
			lastSliceLen: 0,
			wantErr:      ErrSplitSizeOutOfSlice,
		},
		{
			name:         "slice is zero, splitSize not is zero",
			slice:        generateIntSliceByLen(0),
			splitSize:    1,
			wantLen:      0,
			lastSliceLen: 0,
			wantErr:      ErrSplitSizeOutOfSlice,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SplitChunk[int](tc.slice, tc.splitSize)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantLen, len(result))
			assert.Equal(t, tc.lastSliceLen, len(result[len(result)-1]))
		})
	}
}

// cd slice && go test -bench . -benchmem

// BenchmarkSplitChunk 性能测试
func BenchmarkSplitChunk(b *testing.B) {
	// 创建一个长度为 1000000 的切片作为测试数据
	data := make([]int, 1000000)
	// 填充一些数据
	for i := 0; i < len(data); i++ {
		data[i] = i
	}

	// 每次运行基准测试时，调用 SplitChunk 函数
	b.ResetTimer() // 重置计时器，确保不测量初始化时间
	for i := 0; i < b.N; i++ {
		// 测试将数据切割成每块 100 元素的子切片
		_, err := SplitChunk(data, 100)
		if err != nil {
			b.Fatalf("Error splitting chunk: %v", err)
		}
	}
}

func BenchmarkSplitChunkV2(b *testing.B) {
	// 创建一个长度为 1000000 的切片作为测试数据
	data := make([]int, 1000000)
	// 填充一些数据
	for i := 0; i < len(data); i++ {
		data[i] = i
	}

	// 每次运行基准测试时，调用 SplitChunk 函数
	b.ResetTimer() // 重置计时器，确保不测量初始化时间
	for i := 0; i < b.N; i++ {
		// 测试将数据切割成每块 100 元素的子切片
		_, err := SplitChunkV2(data, 100)
		if err != nil {
			b.Fatalf("Error splitting chunk: %v", err)
		}
	}
}

func BenchmarkSplitChunkLow(b *testing.B) {
	// 创建一个长度为 1000000 的切片作为测试数据
	data := make([]int, 1000000)
	// 填充一些数据
	for i := 0; i < len(data); i++ {
		data[i] = i
	}

	// 每次运行基准测试时，调用 SplitChunk 函数
	b.ResetTimer() // 重置计时器，确保不测量初始化时间
	for i := 0; i < b.N; i++ {
		// 测试将数据切割成每块 100 元素的子切片
		_, err := SplitChunkLow(data, 100)
		if err != nil {
			b.Fatalf("Error splitting chunk: %v", err)
		}
	}
}
