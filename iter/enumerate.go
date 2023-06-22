package iter

import (
	"github.com/gpabois/gostd/collection"
	"github.com/gpabois/gostd/option"
)

type Enumeration[T any] collection.Pair[int, T]

func (pair Enumeration[T]) GetFirst() int {
	return pair.First
}

func (pair Enumeration[T]) GetSecond() T {
	return pair.Second
}

type EnumIterator[T any] struct {
	iter  Iterator[T]
	index int
}

func (it *EnumIterator[T]) Next() option.Option[Enumeration[T]] {
	it.index++
	c := it.iter.Next()

	return option.Map(c, func(el T) Enumeration[T] {
		return Enumeration[T]{it.index, el}
	})
}

func Enumerate[T any](iter Iterator[T]) Iterator[Enumeration[T]] {
	return &EnumIterator[T]{
		iter:  iter,
		index: -1,
	}
}
