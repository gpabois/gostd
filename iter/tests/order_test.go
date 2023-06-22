package iter_tests

import (
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/stretchr/testify/assert"
)

func Test_Min_NotEmpty(t *testing.T) {
	it := iter.IterSlice(&[]int{3, 1, 2})
	assert.Equal(t, 1, iter.Min(it).Expect())
}

func Test_Min_Empty(t *testing.T) {
	it := iter.IterSlice(&[]int{})
	assert.True(t, true, iter.Min(it).IsNone())
}

func Test_Max_NotEmpty(t *testing.T) {
	it := iter.IterSlice(&[]int{3, 1, 2})
	assert.Equal(t, 3, iter.Max(it).Expect())
}

func Test_Max_Empty(t *testing.T) {
	it := iter.IterSlice(&[]int{})
	assert.True(t, true, iter.Max(it).IsNone())
}
