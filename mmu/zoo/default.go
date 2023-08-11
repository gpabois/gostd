package zoo

import (
	"github.com/gpabois/gostd/collection/stack"
	"github.com/gpabois/gostd/option"
)

type Cell[T any] struct {
	freed bool
	value T
}

// Default Zoo
type DefaultZoo[T any] struct {
	inWith    bool
	elements  []Cell[T]
	available stack.Stack[uint]
}

// Default zoo
func Default[T any](capacity uint) DefaultZoo[T] {
	return DefaultZoo[T]{
		elements: make([]Cell[T], capacity),
	}
}

// Create a new node
// Don't call within a 'With' call, as the resource could be reallocated.
func (zoo *DefaultZoo[T]) New() uint {
	if zoo.inWith {
		panic("cannot create a new node within a 'With' call.")
	}

	// One is free
	if zoo.available.Last().IsSome() {
		return zoo.available.Pop().Expect()
	}

	nodeId := len(zoo.elements)
	zoo.elements = append(zoo.elements, Cell[T]{})
	return uint(nodeId)
}

// Create a new node
// Don't call within a 'With' call, as the resource could be reallocated.
func (zoo *DefaultZoo[T]) TryNew() option.Option[uint] {
	return option.Some(zoo.New())
}

func (zoo *DefaultZoo[T]) Exists(id uint) bool {
	if int(id) >= len(zoo.elements) {
		return false
	}

	return !zoo.elements[id].freed
}

func (zoo *DefaultZoo[T]) With(nodeId uint, fn func(*T)) {
	if int(nodeId) >= len(zoo.elements) {
		return
	}

	cell := &zoo.elements[nodeId]
	if cell.freed {
		return
	}

	fn(&cell.value)
}

func (zoo *DefaultZoo[T]) Delete(nodeId uint) {
	if int(nodeId) >= len(zoo.elements) {
		return
	}

	cell := &zoo.elements[nodeId]
	if cell.freed {
		return
	}

	cell.freed = true
	zoo.available.Push(nodeId)
}
