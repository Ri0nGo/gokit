package slice

import (
	"github.com/Ri0nGo/gokit/utils"
)

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
