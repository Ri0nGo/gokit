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
