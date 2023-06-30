package norm

import (
	"reflect"
	"strconv"
	"time"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/reflectutil"
	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde/decoder"
)

type Decoder struct {
	m          map[string]any
	dateLayout string
}

func NewDecoder(m map[string]any) *Decoder {
	return &Decoder{m, "YYYY-MM-DDTHH:mm:ss.sssZ"}
}

func (decoder *Decoder) GetCursor() result.Result[any] {
	return result.Success[any](decoder.m)
}

func (decoder *Decoder) IsNull(data any) bool {
	return false
}

func (d *Decoder) DecodeTime(data any, typ reflect.Type) result.Result[reflect.Value] {
	switch v := data.(type) {
	case time.Time:
		return result.Success(reflect.ValueOf(v))
	case string:
		t, err := time.Parse(d.dateLayout, v)
		if err != nil {
			return result.Result[reflect.Value]{}.Failed(err)
		}
		return result.Success(reflect.ValueOf(t))
	default:
		return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(reflect.TypeOf(v), reflectutil.TypeOf[time.Time]()))
	}
}

func (d *Decoder) DecodePrimaryType(data any, typ reflect.Type) result.Result[reflect.Value] {
	switch v := data.(type) {
	case string:
		switch typ.Kind() {
		case reflect.String:
			return result.Success(reflect.ValueOf(v))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return result.Result[reflect.Value]{}.Failed(err)
			}
			return result.Success(reflect.ValueOf(int(i)))
		case reflect.Bool:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return result.Result[reflect.Value]{}.Failed(err)
			}
			return result.Success(reflect.ValueOf(b))

		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return result.Result[reflect.Value]{}.Failed(err)
			}
			return result.Success(reflect.ValueOf(f))
		default:
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(reflect.TypeOf(v), typ))
		}
	case int, int8, int16, int32, int64:
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return result.Success(reflect.ValueOf(v))
		default:
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(reflect.TypeOf(v), typ))
		}
	case bool:
		switch typ.Kind() {
		case reflect.Bool:
			return result.Success(reflect.ValueOf(v))
		default:
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(reflect.TypeOf(v), typ))
		}
	case float32, float64:
		switch typ.Kind() {
		case reflect.Float32, reflect.Float64:
			return result.Success(reflect.ValueOf(v))
		default:
			return result.Result[reflect.Value]{}.Failed(decoder.NewWrongTypeError(reflect.TypeOf(v), typ))
		}
	default:
		return result.Result[reflect.Value]{}.Failed(decoder.NewExpectingPrimaryTypeError(reflect.TypeOf(v)))
	}
}

func (dec *Decoder) IterMap(ast any) result.Result[iter.Iterator[decoder.Element]] {
	switch s := ast.(type) {
	case map[string]any:
		it := iter.Map(iter.IterMap(&s), func(el iter.KV[string, any]) decoder.Element { return el })
		return result.Success(it)
	default:
		return result.Result[iter.Iterator[decoder.Element]]{}.Failed(decoder.NewExpectingMapError(reflect.TypeOf(s)))
	}
}

func (decoder *Decoder) IterSlice(ast any) result.Result[iter.Iterator[any]] {
	return iter.Reflect_Iter(ast)
}
