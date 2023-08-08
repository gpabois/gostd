package result

import (
	"errors"
	"reflect"
)

type IResult interface {
	TypeOf() reflect.Type
	UnwrapError() error
	HasFailed() bool
	IsSuccess() bool
}

type Result[T any] struct {
	inner T
	err   error
}

// Copy error from one result type to another
// Panic if the result is not a failure.
func CopyError[U any, T any](res Result[T]) Result[U] {
	return Failed[U](res.UnwrapError())
}

func FromRaw[T any](val *T, err error) Result[T] {
	if val == nil {
		return Result[T]{}.Failed(err)
	} else {
		return Success(*val)
	}
}

func (res Result[T]) UnwrapRaw() (*T, error) {
	return res.UnwrapValue(), res.UnwrapError()
}

// Unwrap the value, return nil if the result has failed.
func (res Result[T]) UnwrapValue() *T {
	if res.HasFailed() {
		return nil
	}

	return &res.inner
}

func (res Result[T]) Failed(err error) Result[T] {
	return Failed[T](err)
}

func IntoAny[T any](result Result[T]) Result[any] {
	if result.HasFailed() {
		return Failed[any](result.err)
	}

	val := any(result.inner)

	return Success(val)
}

func FromAny[T any](result Result[any]) Result[T] {
	if result.HasFailed() {
		return Failed[T](result.err)
	}

	val, ok := result.inner.(T)

	if !ok {
		return Failed[T](errors.New("wrong type"))
	}

	return Success(val)
}

// Re-enter from Result[any]
func ChainFromAny[T any, U any](inner func(T) Result[U]) func(outer any) Result[any] {
	return func(outer any) Result[any] {
		val, ok := outer.(T)

		if !ok {
			return Failed[any](errors.New("wrong type"))
		}

		return inner(val).ToAny()
	}
}

// Re-enter from Result[any]
func ThenFromAny[T any](inner func(T)) func(outer any) {
	return func(outer any) {
		val, ok := outer.(T)
		if !ok {
			panic(errors.New("wrong type"))
		}

		inner(val)
	}
}

func (result Result[T]) ToAny() Result[any] {
	return Result[any]{
		inner: result.inner,
		err:   result.err,
	}
}

// Result chaining, take a successful result and create another result
// For operations that only require a value transformation without any error, use Map()
func (result Result[T]) Chain(mapper func(val T) Result[T]) Result[T] {
	if result.HasFailed() {
		return Failed[T](result.err)
	} else {
		return mapper(result.Expect())
	}
}

// Flatten a nested result
func Flatten[T any](value Result[Result[T]]) Result[T] {
	if value.HasFailed() {
		return Failed[T](value.UnwrapError())
	} else {
		return value.Expect()
	}
}

func FlatMap[T any, U any](val Result[T], mapper func(val T) Result[U]) Result[U] {
	return Flatten(Map(val, mapper))
}

func Chain[T any, U any](mapper func(val T) Result[U], val Result[T]) Result[U] {
	return FlatMap(val, mapper)
}

// Allow to execute a procedure, while returning the result.
func (res Result[T]) Then(then func(val T)) Result[T] {
	if res.IsSuccess() {
		then(res.Expect())
	}

	return res
}

func Map[T any, U any](val Result[T], mapper func(val T) U) Result[U] {
	if val.HasFailed() {
		return Failed[U](val.err)
	} else {
		return Success(mapper(val.Expect()))
	}
}

func (res Result[T]) HasFailed() bool {
	return res.err != nil
}

func (res Result[T]) IsSuccess() bool {
	return res.err == nil
}

func (res Result[T]) IntoAnyTuple() (any, error) {
	if res.HasFailed() {
		return nil, res.err
	} else {
		return res, nil
	}
}

func (res Result[T]) Expect() T {
	if res.HasFailed() {
		panic(res.err)
	}

	return res.inner
}

func (res Result[T]) UnwrapError() error {
	return res.err
}

func Success[T any](value T) Result[T] {
	return Result[T]{
		inner: value,
		err:   nil,
	}
}

func Failed[T any](err error) Result[T] {
	return Result[T]{
		err: err,
	}
}
