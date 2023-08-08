package result

import (
	"errors"
	"testing"

	"github.com/gpabois/gostd/result"
	"github.com/stretchr/testify/assert"
)

func Test_From_Success(t *testing.T) {
	val := 10
	res := result.FromRaw(&val, nil)
	assert.True(t, res.IsSuccess())
}

func Test_From_Error(t *testing.T) {
	res := result.FromRaw[int](nil, errors.New("error"))
	assert.True(t, res.HasFailed())
}

func Test_Success(t *testing.T) {
	res := result.Success(10)
	assert.True(t, res.IsSuccess())
	assert.Equal(t, 10, res.Expect())

	val, err := res.UnwrapRaw()
	assert.NotNil(t, val)
	assert.Nil(t, err)
}

func Test_Failed(t *testing.T) {
	expectedErr := errors.New("error")
	res := result.Failed[int](expectedErr)
	assert.True(t, res.HasFailed())
	assert.Equal(t, expectedErr, res.UnwrapError())

	// Panics if we expect a value but it has failed.
	assert.Panics(t, func() {
		res.Expect()
	})

	val, err := res.UnwrapRaw()
	assert.Nil(t, val)
	assert.Equal(t, expectedErr, err)

}

func Test_Any_Success(t *testing.T) {
	anyRes := result.Success(10).ToAny()
	res := result.FromAny[int](anyRes)
	assert.True(t, res.IsSuccess())
}

func Test_Any_Failed(t *testing.T) {
	anyRes := result.Success(10).ToAny()
	res := result.FromAny[bool](anyRes)
	assert.True(t, res.HasFailed())
}

func Test_MapResult(t *testing.T) {
	res := result.Map(result.Success(10), func(val int) bool { return val > 9 })
	assert.True(t, res.Expect())
}
