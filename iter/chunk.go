package iter

import "github.com/gpabois/gostd/option"

type Chunk[T any] []T

type chunkIterator[T any] struct {
	chunkSize int
	inner     Iterator[T]
}

func ChunkEvery[T any](it Iterator[T], chunkSize int) Iterator[Chunk[T]] {
	return &chunkIterator[T]{
		chunkSize: chunkSize,
		inner:     it,
	}
}

func (it *chunkIterator[T]) Next() option.Option[Chunk[T]] {
	chunk := Take[Chunk[T]](it.inner, it.chunkSize)

	if len(chunk) == 0 {
		return option.None[Chunk[T]]()
	}

	return option.Some(chunk)
}
