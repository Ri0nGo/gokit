package slice

import (
	"github.com/Ri0nGo/gokit/utils"
)

// SplitChunk[T any] 将slice按照splitSize切割成多段
// 例如：[1,2,3,4,5,6,7,8,9]；splitSize = 4
// 结果：[[1,2,3,4], [5,6,7,8], [9]]
func SplitChunk[T any](slice []T, splitSize int) ([][]T, error) {
	if splitSize <= 0 || splitSize > len(slice) {
		return nil, ErrSplitSizeOutOfSlice
	}

	// splitSize 此时不可能会是0
	quo, _ := utils.RoundUpToQuotient(len(slice), splitSize)

	result := make([][]T, quo)
	for i := 0; i < len(slice); i += splitSize {
		end := i + splitSize
		if end > len(slice) {
			end = len(slice)
		}

		result[i/splitSize] = slice[i:end]
	}
	return result, nil
}

// SplitChunkV2[T any] 将slice按照splitSize切割成多段
// 例如：[1,2,3,4,5,6,7,8,9]；splitSize = 4
// 结果：[[1,2,3,4], [5,6,7,8], [9]]
func SplitChunkV2[T any](slice []T, splitSize int) ([][]T, error) {
	if splitSize <= 0 || splitSize > len(slice) {
		return nil, ErrSplitSizeOutOfSlice
	}

	// splitSize 此时不可能会是0
	quo, _ := utils.RoundUpToQuotient(len(slice), splitSize)

	result := make([][]T, 0, quo)
	for i := 0; i < len(slice); i += splitSize {
		end := i + splitSize
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result, nil
}

func SplitChunkLow[T any](slice []T, splitSize int) ([][]T, error) {
	if splitSize <= 0 || splitSize > len(slice) {
		return nil, ErrSplitSizeOutOfSlice
	}

	result := make([][]T, 0)
	for i := 0; i < len(slice); i += splitSize {
		end := i + splitSize
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[i:end])
	}
	return result, nil
}
