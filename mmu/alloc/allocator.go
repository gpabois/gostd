package alloc

import (
	"unsafe"

	"github.com/gpabois/gostd/option"
)

type Allocator[T any] interface {
	TryNew(count uintptr) option.Option[*T]
	New(count uintptr) *T
	Delete(*T)
}

func Default[T any]() Allocator[T] {
	return DefaultAllocator[T]{}
}

func TryAlloc[T any](alloc Allocator[byte]) option.Option[*T] {
	var t T
	size := unsafe.Sizeof(t)
	return option.Map(alloc.TryNew(size), func(rawPtr *byte) *T {
		return (*T)(unsafe.Pointer(rawPtr))
	})
}

func Alloc[T any](alloc IAllocator[byte]) *T {
	return TryAlloc[T](alloc).Expect()
}
