package alloc

import "github.com/gpabois/gostd/option"

// Used the default GC allocator
type DefaultAllocator[T any] struct{}

func (alloc DefaultAllocator[T]) TryNew() option.Option[*T] {
	var t T
	return option.Some(&t)
}

func (alloc DefaultAllocator[T]) New() *T {
	var t T
	return &t
}

func (alloc DefaultAllocator[T]) Delete(*T) {}
