package norm

import (
	"reflect"
	"time"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/reflectutil"
	cast "github.com/gpabois/gostd/reflectutil/cast"
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
	return cast.Cast(data, typ).IntoResult(decoder.NewWrongTypeError(reflect.TypeOf(data), typ))
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
