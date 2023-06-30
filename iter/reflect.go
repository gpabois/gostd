package iter

import (
	"errors"
	"reflect"

	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/result"
)

type reflectIterator struct {
	val    reflect.Value
	cursor int
}

func (it *reflectIterator) Next() option.Option[any] {
	it.cursor++
	if it.cursor >= it.val.Len() {
		return option.None[any]()
	}
	v := it.val.Index(it.cursor)
	return option.Some(v.Interface())
}

func Reflect_Iter(value any) result.Result[Iterator[any]] {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		return result.Success[Iterator[any]](&reflectIterator{
			val:    v,
			cursor: -1,
		})
	}

	return result.Result[Iterator[any]]{}.Failed(errors.New("not a slice"))
}
