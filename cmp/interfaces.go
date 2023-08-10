package cmp

import (
	"errors"

	"github.com/gpabois/gostd/result"
	"golang.org/x/exp/constraints"
)

const (
	Less = byte(iota)
	Greater
	Equal
)

type Order = byte
type Orderable interface {
	Cmp(right any) result.Result[Order]
}

func cmpOrdered[T constraints.Ordered](left T, right T) result.Result[Order] {
	if left < right {
		return result.Success(Less)
	} else if left > right {
		return result.Success(Greater)
	} else {
		return result.Success(Equal)
	}
}

func anyCmpOrdered(left any, right any) result.Result[Order] {
	switch l := left.(type) {
	case int:
		r := right.(int)
		return cmpOrdered(l, r)
	case int8:
		r := right.(int8)
		return cmpOrdered(l, r)
	case int16:
		r := right.(int16)
		return cmpOrdered(l, r)
	case int32:
		r := right.(int32)
		return cmpOrdered(l, r)
	case int64:
		r := right.(int64)
		return cmpOrdered(l, r)
	case uint:
		r := right.(uint)
		return cmpOrdered(l, r)
	case uint8:
		r := right.(uint8)
		return cmpOrdered(l, r)
	case uint16:
		r := right.(uint16)
		return cmpOrdered(l, r)
	case uint32:
		r := right.(uint32)
		return cmpOrdered(l, r)
	case uint64:
		r := right.(uint64)
		return cmpOrdered(l, r)
	case float32:
		r := right.(float32)
		return cmpOrdered(l, r)
	case string:
		r := right.(string)
		return cmpOrdered(l, r)
	default:
		return result.Failed[Order](errors.New("not comparable"))
	}
}

func Cmp(left any, right any) result.Result[Order] {
	switch l := left.(type) {
	case Orderable:
		return l.Cmp(right)
	default:
		return anyCmpOrdered(left, right)
	}
}

func Max[T any](left T, right T) result.Result[T] {
	return result.Map(Cmp(left, right), func(order Order) T {
		if order == Less {
			return right
		} else {
			return left
		}
	})
}
