package arena

import (
	"unsafe"

	"github.com/gpabois/gostd/option"
)

type cell[T any] struct {
	rc    uint
	value *T
}

func (cell *cell[T]) acquire() Holder[T] {
	cell.rc += 1
	return Holder[T]{
		cell: cell,
	}
}

type Holder[T any] struct {
	released bool
	cell     *cell[T]
}

func (ptr *Holder[T]) Ref() *T {
	return ptr.cell.value
}

// Release the pointer
func (ptr *Holder[T]) Release() {
	if !ptr.released {
		ptr.cell.rc -= 1
	}
}

type SimpleArena[T any] struct {
	elements map[uint]cell[T]
}

func (arena *SimpleArena[T]) New() uint {
	var el T
	ptr := &el

	addr := uint(uintptr(unsafe.Pointer(ptr)))
	arena.elements[addr] = cell[T]{value: ptr}
	return addr
}

func (arena *SimpleArena[T]) At(index uint) option.Option[Holder[T]] {
	cell, ok := arena.elements[index]
	if !ok {
		return option.None[Holder[T]]()
	}
	return option.Some(cell.acquire())
}
