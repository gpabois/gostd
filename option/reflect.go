package option

import (
	"reflect"

	"github.com/gpabois/gostd/reflectutil"
	"github.com/gpabois/gostd/result"
)

func Reflect_Get(val reflect.Value) Option[reflect.Value] {
	opt := val.Interface().(iOption)
	if opt.IsSome() {
		return Some(reflect.ValueOf(opt.Get()))
	} else {
		return None[reflect.Value]()
	}
}

func Reflect_GetInnerType(typ reflect.Type) reflect.Type {
	return reflect.New(typ).Elem().Interface().(iOption).TypeOf()
}

func Reflect_IsOptionType(typ reflect.Type) bool {
	return typ.Implements(reflectutil.TypeOf[iOption]()) && reflect.PtrTo(typ).Implements(reflectutil.TypeOf[iMutableOption]())
}

func Reflect_TrySome(ptrDest reflect.Value, src reflect.Value) result.Result[bool] {
	return ptrDest.Interface().(iMutableOption).TrySome(src.Interface())
}
