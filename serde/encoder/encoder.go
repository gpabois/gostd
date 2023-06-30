package encoder

import (
	"reflect"

	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/result"
)

// Interface for data encoding
// Is a pushdown automata to encode nested data (Map, Array)
// Used mainly by the Encode function, don't try to use it directly.
// TODO : Encode Time
//
//go:generate mockery
type Encoder interface {
	EncodeInt64(value int64) result.Result[bool]
	EncodeFloat64(value float64) result.Result[bool]
	EncodeBool(value bool) result.Result[bool]
	EncodeString(value string) result.Result[bool]
	PushArray() result.Result[bool]
	PushArrayValue() result.Result[bool]
	PushMap() result.Result[bool]
	PushMapKey() result.Result[bool]
	PushMapValue() result.Result[bool]

	Pop() result.Result[bool] // Pop the tos state of the encoder, any Push* must have its Pop, use the with* helpers to ensure it.
}

// Encode any value
func Encode[T any](enc Encoder, value T) {
	encode(enc, reflect.ValueOf(value))
	enc.Pop() // Pop the root state to flush any internal buffer
}

func withArray(encoder Encoder, f func()) {
	encoder.PushArray()
	defer encoder.Pop()
	f()
}

func withArrayValue(encoder Encoder, f func()) {
	encoder.PushArrayValue()
	defer encoder.Pop()
	f()
}

func withMap(encoder Encoder, f func()) {
	encoder.PushMap()
	defer encoder.Pop()
	f()
}

func withMapKey(encoder Encoder, f func()) {
	encoder.PushMapKey()
	defer encoder.Pop()
	f()
}

func withMapValue(encoder Encoder, f func()) {
	encoder.PushMapValue()
	defer encoder.Pop()
	f()
}

func encode(enc Encoder, value reflect.Value) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		enc.EncodeInt64(int64(value.Int()))
	case reflect.String:
		enc.EncodeString(value.String())
	case reflect.Float32, reflect.Float64:
		enc.EncodeFloat64(value.Float())
	case reflect.Bool:
		enc.EncodeBool(value.Bool())
	case reflect.Array, reflect.Slice:
		encodeSlice(enc, value)
	case reflect.Map:
		encodeMap(enc, value)
	case reflect.Struct:
		encodeStruct(enc, value)
	}
}

func encodeStruct(enc Encoder, value reflect.Value) {
	// Encode the struct as a map
	withMap(enc, func() {
		typ := value.Type()

		for i := 0; i < typ.NumField(); i++ {
			field := value.Field(i)
			fieldName := typ.Field(i).Name

			marshalName, ok := typ.Field(i).Tag.Lookup("serde")
			if ok {
				fieldName = marshalName
			}

			if option.Reflect_IsOptionType(field.Type()) {
				optVal := option.Reflect_Get(field)
				if optVal.IsNone() {
					continue
				}

				// Encode element's key
				withMapKey(enc, func() { encode(enc, reflect.ValueOf(fieldName)) })
				// Encode element's value
				withMapValue(enc, func() { encode(enc, optVal.Expect()) })
			} else {
				// Encode element's key
				withMapKey(enc, func() { encode(enc, reflect.ValueOf(fieldName)) })
				// Encode element's value
				withMapValue(enc, func() { encode(enc, field) })
			}
		}
	})
}

func encodeMap(enc Encoder, value reflect.Value) {
	withMap(enc, func() {
		for _, mapKey := range value.MapKeys() {
			// Get the map's value behind the key
			mapValue := value.MapIndex(mapKey)
			// Encode element's key
			withMapKey(enc, func() {
				encode(enc, mapKey)
			})
			// Encode element's value
			withMapValue(enc, func() {
				encode(enc, mapValue)
			})
		}
	})
}
func encodeSlice(enc Encoder, value reflect.Value) {
	withArray(enc, func() {
		for i := 0; i < value.Len(); i++ {
			withArrayValue(enc, func() {
				encode(enc, value.Index(i))
			})
		}
	})
}
