package stack

import "github.com/gpabois/gostd/option"

// LIFO
type IStack[T any] interface {
	Push(value T)
	Pop() option.Option[T]
	Last() option.Option[*T]
}

type Stack[T any] struct {
	inner []T
}

func (stack *Stack[T]) Length() uint {
	return uint(len(stack.inner))
}

func (stack *Stack[T]) Push(value T) {
	stack.inner = append(stack.inner, value)
}

func (stack *Stack[T]) Last() option.Option[*T] {
	if len(stack.inner) == 0 {
		return option.None[*T]()
	}

	return option.Some(&stack.inner[len(stack.inner)-1])
}

func (stack *Stack[T]) Pop() option.Option[T] {
	if len(stack.inner) == 0 {
		return option.None[T]()
	}

	last := stack.inner[len(stack.inner)-1]
	stack.inner = stack.inner[:len(stack.inner)-1]

	return option.Some(last)
}
