package arena

import (
	"unsafe"

	"github.com/gpabois/gostd/option"
)

type simpleCell[T any] struct {
	rc    uint
	value *T
}

func (cell *simpleCell[T]) acquire() simpleHolder[T] {
	cell.rc += 1
	return simpleHolder[T]{
		cell: cell,
	}
}

type simpleHolder[T any] struct {
	released bool
	cell     *simpleCell[T]
}

func (ptr simpleHolder[T]) Ref() *T {
	return ptr.cell.value
}

// Release the pointer
func (ptr simpleHolder[T]) Release() {
	if !ptr.released {
		ptr.cell.rc -= 1
	}
}

func (h simpleHolder[T]) Copy() IHolder[T] {
	return h.cell.acquire()
}

type SimpleArena[T any] struct {
	elements map[uint]simpleCell[T]
}

func NewSimpleArena[T any]() SimpleArena[T] {
	return SimpleArena[T]{
		elements: make(map[uint]simpleCell[T]),
	}
}

func (arena *SimpleArena[T]) New() uint {
	var el T
	ptr := &el

	addr := uint(uintptr(unsafe.Pointer(ptr)))
	arena.elements[addr] = simpleCell[T]{value: ptr}
	return addr
}

func (arena *SimpleArena[T]) At(index uint) option.Option[IHolder[T]] {
	cell, ok := arena.elements[index]
	if !ok {
		return option.None[IHolder[T]]()
	}
	h := cell.acquire()
	return option.Some[IHolder[T]](h)
}

func (arena *SimpleArena[T]) Delete(index uint) {
	cell, ok := arena.elements[index]
	if ok && cell.rc == 0 {
		delete(arena.elements, index)
	}
}
