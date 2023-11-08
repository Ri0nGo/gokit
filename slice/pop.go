package slice

func Pop[T any](src []T, index int) ([]T, T, error) {
	length := len(src)
	if index < 0 || index >= length {
		var zero T
		return nil, zero, ErrIndexOutOfSlice
	}
	popVal := src[index]
	for i := index; i+1 < length; i++ {
		if length-i >= 1 {
			src[i] = src[i+1]
		}
	}
	src = src[:length-1]
	return src, popVal, nil
}
