package slice

type filterFunc[T any] func(elem T) bool
