package cast

import (
	"reflect"
	"strconv"

	"github.com/gpabois/gostd/option"
)

func Cast(from any, typ reflect.Type) option.Option[reflect.Value] {
	switch v := from.(type) {
	case string:
		switch typ.Kind() {
		case reflect.String:
			return option.Some(reflect.ValueOf(v))
		case reflect.Int8:
			i, err := strconv.ParseInt(v, 10, 8)
			if err != nil {
				return option.None[reflect.Value]()
			}
			return option.Some(reflect.ValueOf(int8(i)))
		case reflect.Int16:
			i, err := strconv.ParseInt(v, 10, 16)
			if err != nil {
				return option.None[reflect.Value]()
			}
			return option.Some(reflect.ValueOf(int16(i)))
		case reflect.Int32:
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return option.None[reflect.Value]()
			}
			return option.Some(reflect.ValueOf(int32(i)))
		case reflect.Int, reflect.Int64:
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return option.None[reflect.Value]()
			}
			return option.Some(reflect.ValueOf(int(i)))
		case reflect.Bool:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return option.None[reflect.Value]()
			}
			return option.Some(reflect.ValueOf(b))
		case reflect.Float32:
			f, err := strconv.ParseFloat(v, 32)
			if err != nil {
				return option.None[reflect.Value]()
			}
			return option.Some(reflect.ValueOf(f))
		case reflect.Float64:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return option.None[reflect.Value]()
			}
			return option.Some(reflect.ValueOf(f))
		}
	case int8:
		switch typ.Kind() {
		case reflect.Int:
			return option.Some(reflect.ValueOf(int(v)))
		case reflect.Int8:
			return option.Some(reflect.ValueOf(int8(v)))
		case reflect.Int16:
			return option.Some(reflect.ValueOf(int16(v)))
		case reflect.Int32:
			return option.Some(reflect.ValueOf(int32(v)))
		case reflect.Int64:
			return option.Some(reflect.ValueOf(int64(v)))
		}
	case int16:
		switch typ.Kind() {
		case reflect.Int:
			return option.Some(reflect.ValueOf(int(v)))
		case reflect.Int8:
			return option.Some(reflect.ValueOf(int8(v)))
		case reflect.Int16:
			return option.Some(reflect.ValueOf(int16(v)))
		case reflect.Int32:
			return option.Some(reflect.ValueOf(int32(v)))
		case reflect.Int64:
			return option.Some(reflect.ValueOf(int64(v)))
		}
	case int32:
		switch typ.Kind() {
		case reflect.Int:
			return option.Some(reflect.ValueOf(int(v)))
		case reflect.Int8:
			return option.Some(reflect.ValueOf(int8(v)))
		case reflect.Int16:
			return option.Some(reflect.ValueOf(int16(v)))
		case reflect.Int32:
			return option.Some(reflect.ValueOf(int32(v)))
		case reflect.Int64:
			return option.Some(reflect.ValueOf(int64(v)))
		}
	case int64:
		switch typ.Kind() {
		case reflect.Int:
			return option.Some(reflect.ValueOf(int(v)))
		case reflect.Int8:
			return option.Some(reflect.ValueOf(int8(v)))
		case reflect.Int16:
			return option.Some(reflect.ValueOf(int16(v)))
		case reflect.Int32:
			return option.Some(reflect.ValueOf(int32(v)))
		case reflect.Int64:
			return option.Some(reflect.ValueOf(int64(v)))
		}
	case int:
		switch typ.Kind() {
		case reflect.Int:
			return option.Some(reflect.ValueOf(int(v)))
		case reflect.Int8:
			return option.Some(reflect.ValueOf(int8(v)))
		case reflect.Int16:
			return option.Some(reflect.ValueOf(int16(v)))
		case reflect.Int32:
			return option.Some(reflect.ValueOf(int32(v)))
		case reflect.Int64:
			return option.Some(reflect.ValueOf(int64(v)))
		}
	case bool:
		switch typ.Kind() {
		case reflect.Bool:
			return option.Some(reflect.ValueOf(v))
		}
	case float32:
		switch typ.Kind() {
		case reflect.Float32:
			return option.Some(reflect.ValueOf(v))
		case reflect.Float64:
			return option.Some(reflect.ValueOf(float64(v)))
		case reflect.Int8:
			return option.Some(reflect.ValueOf(int8(v)))
		case reflect.Int16:
			return option.Some(reflect.ValueOf(int16(v)))
		case reflect.Int32:
			return option.Some(reflect.ValueOf(int32(v)))
		case reflect.Int64:
			return option.Some(reflect.ValueOf(int64(v)))
		case reflect.Int:
			return option.Some(reflect.ValueOf(int(v)))
		}
	case float64:
		switch typ.Kind() {
		case reflect.Float32:
			return option.Some(reflect.ValueOf(float32(v)))
		case reflect.Float64:
			return option.Some(reflect.ValueOf(v))
		case reflect.Int8:
			return option.Some(reflect.ValueOf(int8(v)))
		case reflect.Int16:
			return option.Some(reflect.ValueOf(int16(v)))
		case reflect.Int32:
			return option.Some(reflect.ValueOf(int32(v)))
		case reflect.Int64:
			return option.Some(reflect.ValueOf(int64(v)))
		case reflect.Int:
			return option.Some(reflect.ValueOf(int(v)))
		}
	}

	return option.None[reflect.Value]()
}
