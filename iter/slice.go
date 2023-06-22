package iter

import opt "github.com/gpabois/gostd/option"

type SliceIterable[T any, A ~[]T] struct {
	inner A
}

func IntoSliceIterable[T any, A ~[]T](slice A) Iterable[T] {
	return &SliceIterable[T, A]{inner: slice}
}

func (iterable SliceIterable[T, A]) Iter() Iterator[T] {
	return IterSlice(&iterable.inner)
}

func CollectToSlice[S ~[]T, T any](iter Iterator[T]) []T {
	return Reduce(iter, func(slice S, el T) S { return append(slice, el) }, make(S, 0))
}

type SliceIterator[T any, A ~[]T] struct {
	slice  *A
	cursor int
}

func IterSlice[T any, A ~[]T](slice *A) Iterator[T] {
	return &SliceIterator[T, A]{
		slice:  slice,
		cursor: -1,
	}
}

func (iter *SliceIterator[T, A]) Next() opt.Option[T] {
	iter.cursor++

	if iter.cursor >= len(*iter.slice) {
		iter.cursor = len(*iter.slice)
		return opt.None[T]()
	}

	return opt.Some((*iter.slice)[iter.cursor])
}
