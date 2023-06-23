package option

import (
	"errors"
	"reflect"

	"github.com/gpabois/gostd/result"
)

type iMutableOption interface {
	TrySome(value any) result.Result[bool]
}

type iOption interface {
	// Get the inner type of the option
	TypeOf() reflect.Type

	// Return nil if option is none, or the inner value
	Get() any

	IsSome() bool
	IsNone() bool
}

// Swap option/result
func Swap[T any](opt Option[result.Result[T]]) result.Result[Option[T]] {
	if opt.IsNone() {
		return result.Success(None[T]())
	} else {
		res := opt.Expect()
		if res.HasFailed() {
			return result.Failed[Option[T]](res.UnwrapError())
		} else {
			return result.Success(Some(res.Expect()))
		}
	}
}

func IntoResultFunc[T any](err error) func(Option[T]) result.Result[T] {
	return func(opt Option[T]) result.Result[T] { return opt.IntoResult(err) }
}

func (opt Option[T]) IntoResult(err error) result.Result[T] {
	if opt.IsNone() {
		return result.Failed[T](err)
	} else {
		return result.Success(opt.value)
	}
}

func IsOption(value any) bool {
	return Reflect_IsOptionType(reflect.TypeOf(value))
}

type Option[T any] struct {
	isSet bool
	value T
}

// Returns nil, or the value if any
func (opt Option[T]) Get() any {
	if !opt.isSet {
		return nil
	}

	return opt.value
}

func (opt Option[T]) TypeOf() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func (opt *Option[T]) TrySome(val any) result.Result[bool] {
	refVal := reflect.ValueOf(val)

	if refVal.Kind() == reflect.Ptr && refVal.IsNil() {
		opt.isSet = false
		return result.Success(true)
	} else {
		inner, ok := val.(T)
		if !ok {
			return result.Failed[bool](errors.New("cannot cast value to expected type"))
		}
		opt.isSet = true
		opt.value = inner
		return result.Success(true)
	}
}

func (opt Option[T]) IsSome() bool {
	return opt.isSet
}

func (opt Option[T]) IsNone() bool {
	return !opt.isSet
}

func (opt Option[T]) UnwrapOrZero() T {
	var zero T
	if !opt.IsSome() {
		return zero
	}

	return opt.value
}

func (opt Option[T]) UnwrapOr(val func() T) T {
	if !opt.IsSome() {
		return val()
	}

	return opt.value
}

func (opt Option[T]) Expect() T {
	if !opt.IsSome() {
		panic("empty value")
	}

	return opt.value
}

func Map[T any, U any](value Option[T], mapper func(val T) U) Option[U] {
	if value.IsNone() {
		return None[U]()
	} else {
		return Some(mapper(value.value))
	}
}

func Chain[T any, U any](value Option[T], mapper func(val T) Option[U]) Option[U] {
	if value.IsNone() {
		return None[U]()
	} else {
		return mapper(value.value)
	}
}

func Some[T any](value T) Option[T] {
	return Option[T]{isSet: true, value: value}
}

func None[T any]() Option[T] {
	return Option[T]{isSet: false}
}
