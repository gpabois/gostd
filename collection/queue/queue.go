package collection

import "github.com/gpabois/gostd/option"

// FIFO
type IQueue[T any] interface {
	Enqueue(value T)
	Dequeue() option.Option[T]
	First() option.Option[*T]
}

type Queue[T any] struct {
	inner []T
}

func NewQueue[T any](elements ...T) Queue[T] {
	return Queue[T]{
		inner: elements,
	}
}

func (queue *Queue[T]) Enqueue(value T) {
	queue.inner = append(queue.inner, value)
}

func (queue *Queue[T]) Dequeue() option.Option[T] {
	if len(queue.inner) == 0 {
		return option.None[T]()
	}

	first := queue.inner[0]
	queue.inner = queue.inner[1:]

	return option.Some(first)
}

func (queue *Queue[T]) First() option.Option[*T] {
	if len(queue.inner) == 0 {
		return option.None[*T]()
	}

	return option.Some(&queue.inner[0])
}
