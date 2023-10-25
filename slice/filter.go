package slice

// Filter 过滤出执行func后结果为true的slice
func Filter[T any](slice []T, filterFunc filterFunc[T]) []T {
	result := make([]T, 0, len(slice))
	for _, elem := range slice {
		if filterFunc(elem) {
			result = append(result, elem)
		}
	}
	return result
}
