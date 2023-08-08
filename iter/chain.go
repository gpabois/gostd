package iter

import (
	"github.com/gpabois/gostd/collection"
	"github.com/gpabois/gostd/option"
)

type chainedIterator[T any] struct {
	inner collection.Queue[Iterator[T]]
}

func Chain[T any](its ...Iterator[T]) Iterator[T] {
	return &chainedIterator[T]{
		inner: collection.NewQueue(its...),
	}
}

func (it *chainedIterator[E]) Next() option.Option[E] {
	if it.inner.First().IsNone() {
		return option.None[E]()
	}

	next := (*it.inner.First().Expect()).Next()

	if next.IsNone() {
		it.inner.Dequeue()
		return it.Next()
	}

	return next
}
