package decoder_tests

import (
	"reflect"
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/reflectutil"
	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde/decoder"
	"github.com/gpabois/gostd/serde/decoder/mocks"
	"github.com/stretchr/testify/assert"
)

type simple struct {
	Opt      option.Option[bool]
	Val      int `serde:"val"`
	OtherVal int
}

type element struct {
	key string
	val any
}

type elements []element

func (el elements) Iter() iter.Iterator[decoder.Element] {
	return iter.Map(iter.IterSlice(&el), func(el element) decoder.Element { return el })
}

func (el element) Key() string {
	return el.key
}

func (el element) Value() any {
	return el.val
}

func Test_Decoder_Into(t *testing.T) {
	d := mocks.NewDecoder(t)

	d.EXPECT().GetCursor().Return(result.Success[any](0))

	value := simple{
		OtherVal: 10,
	}

	mapElements := elements{{"val", 0}, {"Opt", true}}
	d.EXPECT().IterMap(0).Return(result.Success(mapElements.Iter()))
	d.EXPECT().IterMap(0).Return(result.Success(mapElements.Iter()))

	d.EXPECT().DecodePrimaryType(0, reflectutil.TypeOf[int]()).Return(result.Success(reflect.ValueOf(0)))

	d.EXPECT().IsNull(true).Return(false)
	d.EXPECT().DecodePrimaryType(true, reflectutil.TypeOf[bool]()).Return(result.Success(reflect.ValueOf(true)))

	res := decoder.DecodeInto(d, &value)
	expectedValue := simple{
		Val:      0,
		Opt:      option.Some(true),
		OtherVal: 10,
	}
	assert.True(t, res.IsSuccess(), res.UnwrapError())
	assert.Equal(t, expectedValue, value)
}

func Test_Decoder(t *testing.T) {
	d := mocks.NewDecoder(t)

	d.EXPECT().GetCursor().Return(result.Success[any](0))

	mapElements := elements{{"val", 0}, {"Opt", true}}
	d.EXPECT().IterMap(0).Return(result.Success(mapElements.Iter()))
	d.EXPECT().IterMap(0).Return(result.Success(mapElements.Iter()))

	d.EXPECT().DecodePrimaryType(0, reflectutil.TypeOf[int]()).Return(result.Success(reflect.ValueOf(0)))

	d.EXPECT().IsNull(true).Return(false)
	d.EXPECT().DecodePrimaryType(true, reflectutil.TypeOf[bool]()).Return(result.Success(reflect.ValueOf(true)))

	expectedValue := simple{
		Val:      0,
		Opt:      option.Some(true),
		OtherVal: 5,
	}

	res := decoder.Decode[simple](d)

	assert.True(t, res.IsSuccess(), res.UnwrapError())
	assert.Equal(t, expectedValue, res.Expect())
}
