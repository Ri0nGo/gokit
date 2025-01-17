package slice

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
