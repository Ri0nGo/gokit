package slice

import "errors"

var (
	ErrIndexOutOfSlice     = errors.New("index out of slice length")
	ErrSplitSizeOutOfSlice = errors.New("split size out of slice length")
	ErrSliceEmpty          = errors.New("slice is empty")
)

type filterFunc[T any] func(elem T) bool

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}
