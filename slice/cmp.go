package slice

func Min[T Ordered](slice []T) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, ErrSliceEmpty
	}
	var minVal = slice[0]
	for _, val := range slice[1:] {
		if val < minVal {
			minVal = val
		}
	}
	return minVal, nil
}

func Max[T Ordered](slice []T) (T, error) {
	var zero T
	if len(slice) == 0 {
		return zero, ErrSliceEmpty
	}
	var maxVal = slice[0]
	for _, val := range slice[1:] {
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal, nil
}

func Avg[T Ordered](slice []T) (float64, error) {
	var (
		sum    T
		result float64
	)
	if len(slice) == 0 {
		return result, ErrSliceEmpty
	}
	for _, val := range slice {
		sum += val
	}
	result = float64(sum) / float64(len(slice))
	return result, nil
}

// Contains[T comparable] 判断元素是否在切片中存在
func Contains[T comparable](silce []T, ele T) bool {
	for _, v := range silce {
		if v == ele {
			return true
		}
	}
	return false
}

// SetSlice[T comparable] 去除slice中重复的元素，并保持原来slice的稳定性
// 例如：[1,2,4,2,5]
// 结果：[1,2,4,5]
func SetSlice[T comparable](slice []T) []T {
	var (
		result []T
		tmpMap = make(map[T]struct{})
	)

	for _, v := range slice {
		if _, exists := tmpMap[v]; !exists {
			tmpMap[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
