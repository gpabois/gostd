package decoder

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/result"
)

type Element interface {
	Key() string
	Value() any
}

//go:generate mockery
type Decoder interface {
	// Init the decoder and return data to be decoded along the way
	Init() result.Result[any]
	// Decode time
	DecodeTime(data any, typ reflect.Type) result.Result[reflect.Value]
	// Decode a primary type
	DecodePrimaryType(data any, typ reflect.Type) result.Result[reflect.Value]
	// Iter over encoded slice element
	IterSlice(data any) result.Result[iter.Iterator[any]]
	// Iter over encoded map element (key/value)
	IterMap(data any) result.Result[iter.Iterator[Element]]
}

func searchElement(decoder Decoder, node any, key string) result.Result[option.Option[Element]] {
	res := decoder.IterMap(node)

	if res.HasFailed() {
		return result.Result[option.Option[Element]]{}.Failed(res.UnwrapError())
	}

	return result.Success(iter.Find(
		res.Expect(),
		func(el Element) bool { return el.Key() == key },
	))
}

func decodeSlice(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	valTyp := typ.Elem()

	iterRes := decoder.IterSlice(encoded)

	if iterRes.HasFailed() {
		return result.Result[reflect.Value]{}.Failed(iterRes.UnwrapError())
	}

	res := iter.Result_FromIter[[]reflect.Value](
		iter.Map(
			iterRes.Expect(),
			func(encoded any) result.Result[reflect.Value] {
				return decode(decoder, encoded, valTyp)
			},
		),
	)

	arr := reflect.New(typ)
	for _, el := range res.Expect() {
		arr.Elem().Set(reflect.Append(arr.Elem(), el))
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

func decodeMap(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	val := reflect.New(typ)

	res := decodeMapElements(decoder, encoded, typ)

	if res.HasFailed() {
		return result.Result[reflect.Value]{}.Failed(res.UnwrapError())
	}

	for _, el := range res.Expect() {
		val.SetMapIndex(reflect.ValueOf(el.Key), el.Value)
	}

	return result.Success(val.Elem())
}

func decodeStruct(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	val := reflect.New(typ)
	for i := 0; i < typ.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		marshalName, ok := typ.Field(i).Tag.Lookup("serde")
		if ok {
			fieldName = marshalName
		}

		cOptRes := searchElement(decoder, encoded, fieldName)
		if cOptRes.HasFailed() {
			return result.Result[reflect.Value]{}.Failed(cOptRes.UnwrapError())
		}
		cOpt := cOptRes.Expect()

		if cOpt.IsNone() || !field.IsValid() {
			continue
		}

		// Decode option
		if field.CanAddr() && option.IsMutableOption(field.Addr().Interface()) {
			innerType := field.Interface().(option.IOption).TypeOf()
			res := decode(decoder, cOpt.Expect(), innerType)
			if res.HasFailed() {
				return result.Failed[reflect.Value](res.UnwrapError())
			}
			resSet := field.Addr().Interface().(option.IMutableOption).TrySet(res.Expect())
			if resSet.HasFailed() {
				return result.Failed[reflect.Value](res.UnwrapError())
			}
		} else { // Decode normally
			res := decode(decoder, cOpt.Expect(), field.Type())
			if res.HasFailed() {
				return result.Failed[reflect.Value](res.UnwrapError())
			}
			field.Set(res.Expect())
		}
	}

	return result.Success(val.Elem())
}

// Try to decode time
func decodeTime(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	return decoder.DecodeTime(encoded, typ)
}

func decode(decoder Decoder, encoded any, typ reflect.Type) result.Result[reflect.Value] {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.Float32, reflect.Float64, reflect.String:
		return decoder.DecodePrimaryType(encoded, typ)
	case reflect.Array, reflect.Slice:
		return decodeSlice(decoder, encoded, typ)
	case reflect.Map:
		return decodeMap(decoder, encoded, typ)
	case reflect.Struct:
		// Decode time
		var t time.Time
		if typ == reflect.TypeOf(t) {
			return decodeTime(decoder, encoded, typ)
		}
		// Decode as a regular struct
		return decodeStruct(decoder, encoded, typ)
	default:
		return result.Result[reflect.Value]{}.Failed(errors.New(fmt.Sprintf("type %v cannot be denormalised", typ.Kind())))
	}
}

func Decode[T any](decoder Decoder) result.Result[T] {
	var v T
	initVal := decoder.Init()
	if initVal.HasFailed() {
		return result.Result[T]{}.Failed(initVal.UnwrapError())
	}
	resVal := decode(decoder, initVal.Expect(), reflect.TypeOf(v))
	if resVal.HasFailed() {
		return result.Result[T]{}.Failed(resVal.UnwrapError())
	}
	return result.Success(resVal.Expect().Interface().(T))
}
