package option_tests

import (
	"testing"

	"github.com/gpabois/gostd/option"
	"github.com/stretchr/testify/assert"
)

func Test_Option_Some(t *testing.T) {
	opt := option.Some(10)

	assert.True(t, opt.IsSome())
	assert.Equal(t, 10, opt.Expect())
}

func Test_Option_None(t *testing.T) {
	opt := option.None[int]()

	assert.True(t, opt.IsNone())

	assert.Panics(t, func() {
		opt.Expect()
	})
}

func Test_Option_Map_Some(t *testing.T) {
	opt := option.Map(option.Some(10), func(e int) bool { return e >= 10 })
	assert.True(t, opt.IsSome())
	assert.Equal(t, true, opt.Expect())
}

func Test_Option_Map_None(t *testing.T) {
	opt := option.Map(option.None[int](), func(e int) bool { return e >= 10 })
	assert.True(t, opt.IsNone())
}

func Test_Option_Chain_Some(t *testing.T) {
	opt := option.Chain(option.Some(10), func(val int) option.Option[bool] {
		if val >= 10 {
			return option.None[bool]()
		}

		return option.Some(true)
	})
	assert.True(t, opt.IsNone())
}

func Test_Option_Chain_None(t *testing.T) {
	opt := option.Chain(option.Some(10), func(val int) option.Option[bool] {
		if val >= 10 {
			return option.Some(true)
		}

		return option.None[bool]()
	})
	assert.True(t, opt.IsSome())
	assert.Equal(t, true, opt.Expect())
}
