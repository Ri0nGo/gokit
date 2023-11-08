package slice

import "errors"

// AddLow[T any] 性能较差，可以执行benchmark进行对比
func InsertLow[T any](slice []T, val T, index int) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return nil, errors.New("index out of slice length")
	}

	result := make([]T, 0, len(slice)+1)
	result = append(result, slice[:index]...)
	result = append(result, val)
	result = append(result, slice[index:]...)
	return result, nil
}

func Insert[T any](src []T, element T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || index >= length {
		return nil, errors.New("index out of slice length")
	}

	//先将src扩展一个元素
	var zeroValue T
	src = append(src, zeroValue)
	for i := len(src) - 1; i > index; i-- {
		if i-1 >= 0 {
			src[i] = src[i-1]
		}
	}
	src[index] = element
	return src, nil
}
