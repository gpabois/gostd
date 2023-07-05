package decoder

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/reflectutil"
	"github.com/gpabois/gostd/result"
)

type Element interface {
	Key() string
	Value() any
}

//go:generate mockery
type Decoder interface {
	// Return cursor to data to be decoded
	GetCursor() result.Result[any]
	// Check if null data (for option)
	IsNull(data any) bool
	// Decode time
	DecodeTime(data any, typ reflect.Type) result.Result[reflect.Value]
	// Decode a primary type
	DecodePrimaryType(data any, typ reflect.Type) result.Result[reflect.Value]
	// Iter over encoded slice element
	IterSlice(data any) result.Result[iter.Iterator[any]]
	// Iter over encoded map element (key/value)
	IterMap(data any) result.Result[iter.Iterator[Element]]
}

func searchElement(node any, key string, elements iter.Iterator[Element]) result.Result[option.Option[Element]] {
	return result.Success(iter.Find(elements,
		func(el Element) bool {
			return el.Key() == key
		},
	))
}

func decodeSliceWithPtr(decoder Decoder, encoded any, arrPtr reflect.Value) result.Result[bool] {
	// Get the inner type of the slice
	elType := arrPtr.Type().Elem().Elem()

	//
	iterRes := decoder.IterSlice(encoded)
	if iterRes.HasFailed() {
		return result.Result[bool]{}.Failed(iterRes.UnwrapError())
	}
	// Decode each element
	res := iter.Result_FromIter[[]reflect.Value](
		iter.Map(
			iterRes.Expect(),
			func(encoded any) result.Result[reflect.Value] {
				return decode(decoder, encoded, elType)
			},
		),
	)

	if res.HasFailed() {
		return result.Result[bool]{}.Failed(res.UnwrapError())
	}

	for _, el := range res.Expect() {
		arrPtr.Elem().Set(reflect.Append(arrPtr.Elem(), el))
	}

	return result.Success(true)
}

func decodeSlice(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	arr := reflect.New(typ)
	if res := decodeSliceWithPtr(decoder, encoded, arr); res.HasFailed() {
		return result.Result[reflect.Value]{}.Failed(res.UnwrapError())
	}
	return result.Success(arr.Elem())
}

type reflectElement struct {
	Key   string
	Value reflect.Value
}

func decodeMapElements(decoder Decoder, encoded any, typ reflect.Type) result.Result[[]reflectElement] {
	iterRes := decoder.IterMap(encoded)

	if iterRes.HasFailed() {
		return result.Result[[]reflectElement]{}.Failed(iterRes.UnwrapError())
	}

	return iter.Result_FromIter[[]reflectElement](
		iter.Map(
			iterRes.Expect(),
			func(element Element) result.Result[reflectElement] {
				return result.Map(
					decode(decoder, element.Value(), typ.Elem()),
					func(decoded reflect.Value) reflectElement {
						return reflectElement{
							Key:   element.Key(),
							Value: decoded,
						}
					},
				)
			},
		),
	)
}

func decodeMapWithPtr(decoder Decoder, encoded any, mapPtr reflect.Value) result.Result[bool] {
	typ := mapPtr.Type().Elem()

	res := decodeMapElements(decoder, encoded, typ)
	if res.HasFailed() {
		return result.Result[bool]{}.Failed(res.UnwrapError())
	}

	for _, el := range res.Expect() {
		mapPtr.SetMapIndex(reflect.ValueOf(el.Key), el.Value)
	}

	return result.Success(true)
}

func decodeMap(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	mapPtr := reflect.New(typ)

	res := decodeMapWithPtr(decoder, encoded, mapPtr)

	if res.HasFailed() {
		return result.Result[reflect.Value]{}.Failed(res.UnwrapError())
	}

	return result.Success(mapPtr.Elem())
}

func decodeStructWithPtr(decoder Decoder, encoded any, structPtr reflect.Value) result.Result[bool] {
	typ := structPtr.Type().Elem()

	elementsRes := decoder.IterMap(encoded)
	if elementsRes.HasFailed() {
		return result.Result[bool]{}.Failed(elementsRes.UnwrapError())
	}
	elements := iter.CollectToSlice[[]Element](elementsRes.Expect())

	for i := 0; i < typ.NumField(); i++ {
		field := structPtr.Elem().Field(i).Addr()
		fieldName := typ.Field(i).Name

		marshalName, ok := typ.Field(i).Tag.Lookup("serde")

		if ok {
			fieldName = marshalName
		}

		cOptRes := searchElement(encoded, fieldName, iter.IterSlice(&elements))

		if cOptRes.HasFailed() {
			return result.Result[bool]{}.Failed(cOptRes.UnwrapError())
		}

		cOpt := cOptRes.Expect()
		if cOpt.IsNone() {
			continue
		}

		res := decodeWithPtr(decoder, cOpt.Expect().Value(), field)

		if res.HasFailed() {
			return result.Failed[bool](res.UnwrapError())
		}
	}

	return result.Success(true)
}

func decodeStruct(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	val := reflect.New(typ)
	decodeStructWithPtr(decoder, encoded, val)
	return result.Success(val.Elem())
}

// Decode optional value
func decodeOption(decoder Decoder, encoded any, optType reflect.Type) result.Result[reflect.Value] {
	ptrOpt := reflect.New(optType)

	if decoder.IsNull(encoded) {
		return result.Success(ptrOpt.Elem())
	}

	innerType := option.Reflect_GetInnerType(optType)
	decRes := decode(decoder, encoded, innerType)

	if decRes.HasFailed() {
		return result.Result[reflect.Value]{}.Failed(decRes.UnwrapError())
	}

	setRes := option.Reflect_TrySome(ptrOpt, decRes.Expect())

	if setRes.HasFailed() {
		return result.Result[reflect.Value]{}.Failed(setRes.UnwrapError())
	}

	return result.Success(ptrOpt.Elem())
}

// Try to decode time
func decodeTime(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	return decoder.DecodeTime(encoded, typ)
}

func decodeWithPtr(decoder Decoder, encoded any, ptr reflect.Value) result.Result[bool] {
	typ := ptr.Type().Elem()
	if typ == reflect.TypeOf(encoded) {
		ptr.Elem().Set(reflect.ValueOf(ptr))
		return result.Success(true)
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.Float32, reflect.Float64, reflect.String:
		res := decoder.DecodePrimaryType(encoded, typ)
		if res.HasFailed() {
			return result.Result[bool]{}.Failed(res.UnwrapError())
		}
		ptr.Elem().Set(res.Expect())
		return result.Success(true)
	case reflect.Array, reflect.Slice:
		return decodeSliceWithPtr(decoder, encoded, ptr)
	case reflect.Map:
		return decodeMapWithPtr(decoder, encoded, ptr)
	case reflect.Struct:
		if typ == reflectutil.TypeOf[time.Time]() { // Decode time
			res := decodeTime(decoder, encoded, typ)
			if res.HasFailed() {
				return result.Result[bool]{}.Failed(res.UnwrapError())
			}
			ptr.Elem().Set(res.Expect())
			return result.Success(true)
		} else if option.Reflect_IsOptionType(typ) { // Decode option
			res := decodeOption(decoder, encoded, typ)
			if res.HasFailed() {
				return result.Result[bool]{}.Failed(res.UnwrapError())
			}
			ptr.Elem().Set(res.Expect())
			return result.Success(true)
		}
		// Decode as a regular struct
		return decodeStructWithPtr(decoder, encoded, ptr)
	default:
		return result.Result[bool]{}.Failed(errors.New(fmt.Sprintf("type %v cannot be denormalised", typ.Kind())))
	}
}

func decode(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	if typ == reflect.TypeOf(encoded) {
		return result.Success(reflect.ValueOf(encoded))
	}
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.Float32, reflect.Float64, reflect.String:
		return decoder.DecodePrimaryType(encoded, typ)
	case reflect.Array, reflect.Slice:
		return decodeSlice(decoder, encoded, typ)
	case reflect.Map:
		return decodeMap(decoder, encoded, typ)
	case reflect.Struct:
		if typ == reflectutil.TypeOf[time.Time]() { // Decode time
			return decodeTime(decoder, encoded, typ)
		} else if option.Reflect_IsOptionType(typ) { // Decode option
			return decodeOption(decoder, encoded, typ)
		}
		// Decode as a regular struct
		return decodeStruct(decoder, encoded, typ)
	default:
		return result.Result[reflect.Value]{}.Failed(errors.New(fmt.Sprintf("type %v cannot be denormalised", typ.Kind())))
	}
}

func Reflect_Decode(decoder Decoder, typ reflect.Type) result.Result[any] {
	cursorRes := decoder.GetCursor()
	if cursorRes.HasFailed() {
		return result.Result[any]{}.Failed(cursorRes.UnwrapError())
	}

	ptr := reflect.New(typ)
	res := decodeWithPtr(decoder, cursorRes.Expect(), ptr)

	if res.HasFailed() {
		return result.Result[any]{}.Failed(res.UnwrapError())
	}

	return result.Success(ptr.Elem().Interface())
}

func Reflect_DecodeInto(decoder Decoder, value any) result.Result[bool] {
	cursorRes := decoder.GetCursor()
	if cursorRes.HasFailed() {
		return result.Result[bool]{}.Failed(cursorRes.UnwrapError())
	}

	ptr := reflect.ValueOf(value)
	return decodeWithPtr(decoder, cursorRes.Expect(), ptr)
}

func DecodeInto[T any](decoder Decoder, ptr *T) result.Result[bool] {
	cursorRes := decoder.GetCursor()
	if cursorRes.HasFailed() {
		return result.Result[bool]{}.Failed(cursorRes.UnwrapError())
	}
	refPtr := reflect.ValueOf(ptr)
	return decodeWithPtr(decoder, cursorRes.Expect(), refPtr)
}

func Decode[T any](decoder Decoder) result.Result[T] {
	var v T
	cursorRes := decoder.GetCursor()
	if cursorRes.HasFailed() {
		return result.Result[T]{}.Failed(cursorRes.UnwrapError())
	}
	resVal := decode(decoder, cursorRes.Expect(), reflect.TypeOf(v))
	if resVal.HasFailed() {
		return result.Result[T]{}.Failed(resVal.UnwrapError())
	}
	return result.Success(resVal.Expect().Interface().(T))
}
