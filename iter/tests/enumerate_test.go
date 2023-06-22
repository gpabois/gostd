package iter_tests

import (
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/stretchr/testify/assert"
)

func Test_Enumerate(t *testing.T) {
	expectedValue := []iter.Enumeration[int]{
		{First: 0, Second: 1},
		{First: 1, Second: 2},
		{First: 2, Second: 3},
	}

	it := iter.IterSlice(&[]int{1, 2, 3})
	value := iter.CollectToSlice[[]iter.Enumeration[int]](iter.Enumerate(it))

	assert.Equal(t, expectedValue, value)
}
