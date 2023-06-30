package json

import (
	"errors"
	"io"
	"reflect"
	"time"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde/decoder"
)

type Decoder struct {
	parser     Parser
	dateLayout string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		parser:     *NewParser(r),
		dateLayout: "YYYY-MM-DDTHH:mm:ss.sssZ", // Default is ISO String
	}
}

func (decoder *Decoder) GetCursor() result.Result[any] {
	return decoder.parser.Parse().ToAny()
}

func (decoder *Decoder) IsNull(data any) bool {
	switch val := data.(type) {
	case Value:
		return val.IsNull()
	default:
		return false
	}
}

func (decoder *Decoder) DecodeTime(data any, typ reflect.Type) result.Result[reflect.Value] {
	switch val := data.(type) {
	case Value:
		if val.IsInteger() {
			return decoder.DecodeTime(val.ExpectInteger(), typ)
		} else if val.IsString() {
			return decoder.DecodeTime(val.ExpectString(), typ)
		}
	// Parse integer as a unix epoch in seconds
	case int:
		return result.Success(reflect.ValueOf(time.Unix(int64(val), 0)))
	// Try to parse the date
	case string:
		t, err := time.Parse(decoder.dateLayout, val)
		if err != nil {
			return result.Result[reflect.Value]{}.Failed(err)
		}
		return result.Success(reflect.ValueOf(t))
	}

	return result.Result[reflect.Value]{}.Failed(errors.New("cannot decode time"))
}

func (decoder *Decoder) DecodePrimaryType(data any, typ reflect.Type) result.Result[reflect.Value] {
	return decodePrimaryTypes(data, typ)
}

func (dec *Decoder) IterMap(ast any) result.Result[iter.Iterator[decoder.Element]] {
	switch node := ast.(type) {
	case Json:
		if !node.IsDocument() {
			return result.Result[iter.Iterator[decoder.Element]]{}.Failed(decoder.NewExpectingMapError(reflect.TypeOf(node)))
		}
		return dec.IterMap(node.ExpectDocument())
	case Value:
		if !node.IsDocument() {
			return result.Result[iter.Iterator[decoder.Element]]{}.Failed(decoder.NewExpectingMapError(reflect.TypeOf(node)))
		}
		return dec.IterMap(node.ExpectArray())
	case Document:
		return result.Success(iter.Map(
			iter.IterSlice(&node.Pairs),
			func(pair Element) decoder.Element {
				return pair
			},
		))
	default:
		return result.Result[iter.Iterator[decoder.Element]]{}.Failed(decoder.NewExpectingMapError(reflect.TypeOf(node)))
	}
}

func (d *Decoder) IterSlice(ast any) result.Result[iter.Iterator[any]] {
	switch node := ast.(type) {
	case Json:
		if !node.IsArray() {
			return result.Result[iter.Iterator[any]]{}.Failed(decoder.NewExpectingSliceError(reflect.TypeOf(node)))
		}
		return d.IterSlice(node.ExpectArray())
	case Value:
		if !node.IsArray() {
			return result.Result[iter.Iterator[any]]{}.Failed(decoder.NewExpectingSliceError(reflect.TypeOf(node)))
		}
		return d.IterSlice(node.ExpectArray())
	case Array:
		return result.Success(iter.Map(iter.IterSlice(&node.Elements), func(el Value) any { return any(el) }))
	default:
		return result.Result[iter.Iterator[any]]{}.Failed(decoder.NewExpectingSliceError(reflect.TypeOf(node)))
	}
}

func decodePrimaryTypes(ast any, typ reflect.Type) result.Result[reflect.Value] {
	val, ok := ast.(Value)

	if !ok {
		return result.Result[reflect.Value]{}.Failed(errors.New("not a value"))
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !val.IsInteger() {
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(typ, reflect.TypeOf(val)))
		}

		return result.Success(reflect.ValueOf(val.ExpectInteger()))
	case reflect.Float32, reflect.Float64:
		if !val.IsFloat() {
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(typ, reflect.TypeOf(val)))
		}
		return result.Success(reflect.ValueOf(val.ExpectFloat()))
	case reflect.Bool:
		if !val.IsBoolean() {
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(typ, reflect.TypeOf(val)))
		}
		return result.Success(reflect.ValueOf(val.ExpectBoolean()))
	case reflect.String:
		if !val.IsString() {
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(typ, reflect.TypeOf(val)))
		}
		return result.Success(reflect.ValueOf(val.ExpectString()))

	default:
		return result.Result[reflect.Value]{}.Failed(decoder.NewExpectingPrimaryTypeError(reflect.TypeOf(val)))
	}
}
