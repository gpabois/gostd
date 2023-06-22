package iter_tests

import (
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/stretchr/testify/assert"
)

func Test_Group(t *testing.T) {
	it := iter.IterSlice(&[]int{1, 2, 3, 4})
	g := iter.Group[[]int](it, func(el int) bool { return el%2 == 0 })

	assert.Equal(t, g[true], []int{2, 4})
	assert.Equal(t, g[false], []int{1, 3})
}
