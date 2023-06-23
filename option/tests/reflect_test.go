package option_tests

import (
	"reflect"
	"testing"

	"github.com/gpabois/gostd/option"
	"github.com/stretchr/testify/assert"
)

func Test_Reflect_IsOptionType(t *testing.T) {
	opt := option.Some(true)

	assert.False(t, option.Reflect_IsOptionType(reflect.TypeOf(10)))
	assert.True(t, option.Reflect_IsOptionType(reflect.TypeOf(opt)))
}

func Test_Reflect_GetInnerType(t *testing.T) {
	opt := option.Some(true)

	assert.Equal(t, reflect.TypeOf(true), option.Reflect_GetInnerType(reflect.TypeOf(opt)))
}

func Test_Reflect_Get(t *testing.T) {
	opt := option.Some(true)

	assert.Equal(t, reflect.ValueOf(true), option.Reflect_Get(reflect.ValueOf(opt)).Expect())
}

func Test_Reflect_TrySome(t *testing.T) {
	opt := option.Some(false)
	ptrOpt := &opt

	res := option.Reflect_TrySome(reflect.ValueOf(ptrOpt), reflect.ValueOf(true))
	assert.True(t, res.IsSuccess(), res.UnwrapError())

	assert.Equal(t, true, ptrOpt.Expect())
}
