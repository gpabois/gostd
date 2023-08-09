package alloc

import (
	"github.com/gpabois/gostd/result"
)

type RawAllocator = Allocator[byte]

type Allocator[T any] interface {
	Alloc(count uint) result.Result[*T]
	Free(value *T)
}
