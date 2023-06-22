package iter_tests

import (
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/stretchr/testify/assert"
)

func Test_Range(t *testing.T) {
	expectedSlice := []int{1, 2, 3, 4}
	slice := iter.CollectToSlice[[]int](iter.Range(1, 4, 1))
	assert.Equal(t, expectedSlice, slice)
}
