package iter

import (
	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/option"
)

type Iterable[T any] interface {
	Iter() Iterator[T]
}

type Iterator[T any] interface {
	Next() option.Option[T]
}

type FromIterator[T any] interface {
	FromIter(iter Iterator[T])
}

func ForEach[T any](iter Iterator[T], fn func(T)) {
	for el := iter.Next(); el.IsSome(); el = iter.Next() {
		fn(el.Expect())
	}
}

func First[T any](iter Iterator[T]) option.Option[T] {
	for el := iter.Next(); el.IsSome(); {
		return option.Some(el.Expect())
	}
	return option.None[T]()
}

func Filter[T any](iter Iterator[T], filter func(T) bool) Iterator[T] {
	return &FilteredIterator[T]{
		inner:  iter,
		filter: filter,
	}
}

func Find[T any](iter Iterator[T], filter func(T) bool) option.Option[T] {
	return Filter(iter, filter).Next()
}

// Take n element in the iterator (or less)
func Take[S ~[]T, T any](it Iterator[T], n int) S {
	var slc S

	for c, ctr := it.Next(), 0; c.IsSome() && ctr <= n-1; c, ctr = it.Next(), ctr+1 {
		slc = append(slc, c.Expect())
		if ctr == n-1 {
			break
		}
	}

	return slc
}

func Reduce[R any, T any](iter Iterator[T], reducer func(R, T) R, init R) R {
	agg := init
	for c := iter.Next(); c.IsSome(); c = iter.Next() {
		v := c.Expect()
		agg = reducer(agg, v)
	}
	return agg
}

func Any(iter Iterator[bool]) bool {
	return Find(iter, ops.IsTrue).UnwrapOr(func() bool { return false })
}

func All(iter Iterator[bool]) bool {
	return Find(iter, ops.IsFalse).UnwrapOr(func() bool { return true })
}

func Map[T any, R any](inner Iterator[T], mapper func(T) R) Iterator[R] {
	return MappedIterator[T, R]{
		mapper,
		inner,
	}
}

func Collect[C, T any](fromIter func(iter Iterator[T]) C, iter Iterator[T]) C {
	return fromIter(iter)
}

type MappedIterator[T any, R any] struct {
	mapper func(T) R
	inner  Iterator[T]
}

func (iter MappedIterator[T, R]) Next() option.Option[R] {
	el := iter.inner.Next()
	if el.IsNone() {
		return option.None[R]()
	}
	val := el.Expect()
	return option.Some(iter.mapper(val))
}

type FilteredIterator[T any] struct {
	filter func(T) bool
	inner  Iterator[T]
}

func (iter *FilteredIterator[T]) Next() option.Option[T] {
	for el := iter.inner.Next(); el.IsSome(); el = iter.inner.Next() {
		value := el.Expect()
		if iter.filter(value) {
			return el
		}
	}
	return option.None[T]()
}
