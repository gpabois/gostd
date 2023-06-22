package bson

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
	parser     *Parser
	dateLayout string
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{parser: NewParser(r), dateLayout: "YYYY-MM-DDTHH:mm:ss.sssZ"}
}

func (d *Decoder) Init() result.Result[any] {
	return d.parser.Parse().ToAny()
}

func (d *Decoder) DecodeTime(ast any, typ reflect.Type) result.Result[reflect.Value] {
	switch node := ast.(type) {
	case Value:
		if node.IsTime() {
			return result.Success(reflect.ValueOf(node.ExpectTime()))
		} else if node.IsInteger() {
			return d.DecodeTime(node.ExpectInteger(), typ)
		} else if node.IsString() {
			return d.DecodeTime(node.ExpectString(), typ)
		}
	// Parse integer as a unix epoch in seconds
	case int:
		return result.Success(reflect.ValueOf(time.Unix(int64(node), 0)))
	// Try to parse the date
	case string:
		t, err := time.Parse(d.dateLayout, node)
		if err != nil {
			return result.Result[reflect.Value]{}.Failed(err)
		}
		return result.Success(reflect.ValueOf(t))

	}
	return result.Result[reflect.Value]{}.Failed(errors.New("cannot decode time"))
}

func (d *Decoder) DecodePrimaryType(ast any, typ reflect.Type) result.Result[reflect.Value] {
	return decodePrimaryTypes(ast, typ)
}

func (d *Decoder) IterMap(ast any) result.Result[iter.Iterator[decoder.Element]] {
	switch node := ast.(type) {
	case Value:
		if node.IsDocument() {
			return d.IterMap(node.ExpectDocument())
		}
	case Document:
		return result.Success(
			iter.Map(iter.IterSlice(&node.Elements), func(el Element) decoder.Element { return el }),
		)
	}

	return result.Result[iter.Iterator[decoder.Element]]{}.Failed(errors.New("not a map"))
}

func decodePrimaryTypes(ast any, typ reflect.Type) result.Result[reflect.Value] {
	val, ok := ast.(Value)

	if !ok {
		return result.Result[reflect.Value]{}.Failed(errors.New("not a value"))
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !val.IsInteger() {
			return result.Result[reflect.Value]{}.Failed(errors.New("not an integer"))
		}

		return result.Success(reflect.ValueOf(val.ExpectInteger()))
	case reflect.Float32, reflect.Float64:
		if !val.IsFloat() {
			return result.Result[reflect.Value]{}.Failed(errors.New("not a float"))
		}
		return result.Success(reflect.ValueOf(val.ExpectFloat()))
	case reflect.Bool:
		if !val.IsBoolean() {
			return result.Result[reflect.Value]{}.Failed(errors.New("not a boolean"))
		}
		return result.Success(reflect.ValueOf(val.ExpectBoolean()))
	case reflect.String:
		if !val.IsString() {
			return result.Result[reflect.Value]{}.Failed(errors.New("not a string"))
		}
		return result.Success(reflect.ValueOf(val.ExpectString()))
	default:
		return result.Result[reflect.Value]{}.Failed(errors.New("not a primary type"))
	}
}
