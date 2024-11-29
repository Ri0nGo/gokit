package slice

import "errors"

var (
	ErrIndexOutOfSlice     = errors.New("index out of slice length")
	ErrSplitSizeOutOfSlice = errors.New("split size out of slice length")
)

type filterFunc[T any] func(elem T) bool
