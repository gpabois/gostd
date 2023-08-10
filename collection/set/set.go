package collection

import (
	"github.com/gpabois/gostd/graph/tree"
	"github.com/gpabois/gostd/graph/tree/rb"
	"github.com/gpabois/gostd/iter"
)

// Unordered collection which avoid any duplicates
type ISet[T any] interface {
	Add(value T)
	Discard(value T)
	Iter() iter.Iterator[T]
}

// Collection which avoid any duplicates
type Set[T any] struct {
	inner tree.SearchTree[T]
}

func NewSet[T any]() Set[T] {
	tree := rb.NewTree[T]()
	return Set[T]{
		inner: &tree,
	}
}

func (set *Set[T]) Add(value T) {
	if set.inner.Search(value).IsSome() {
		set.inner.Remove(value)
	}
	set.inner.Insert(value)
}

func (set *Set[T]) Discard(value T) {
	set.inner.Remove(value)
}
