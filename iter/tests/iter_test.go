package iter_tests

import (
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/option"
	"github.com/stretchr/testify/assert"
)

func Test_IterSlice_Success(t *testing.T) {
	expectedSlice := []int{1, 2, 3}
	it := iter.IterSlice(&expectedSlice)
	slice := iter.CollectToSlice[[]int](it)
	assert.Equal(t, expectedSlice, slice)
}

func Test_Map_Success(t *testing.T) {
	expectedSlice := []int{2, 3, 4}
	it := iter.Map(iter.IterSlice(&[]int{1, 2, 3}), func(el int) int { return el + 1 })
	slice := iter.CollectToSlice[[]int](it)
	assert.Equal(t, expectedSlice, slice)
}

func Test_Filter_Success(t *testing.T) {
	expectedSlice := []int{2}
	it := iter.Filter(iter.IterSlice(&[]int{1, 2, 3}), func(el int) bool { return el == 2 })
	slice := iter.CollectToSlice[[]int](it)
	assert.Equal(t, expectedSlice, slice)
}

func Test_Find_Success(t *testing.T) {
	expectedValue := option.Some(2)
	value := iter.Find(iter.IterSlice(&[]int{1, 2, 3}), func(el int) bool { return el == 2 })
	assert.Equal(t, expectedValue, value)
}

func Test_Reduce_Success(t *testing.T) {
	expected := 6
	value := iter.Reduce(iter.IterSlice(&[]int{1, 2, 3}), ops.Add2[int], 0)
	assert.Equal(t, expected, value)
}

func Test_All_True(t *testing.T) {
	slice := []bool{true, true, true}
	value := iter.All(iter.IterSlice(&slice))
	assert.Equal(t, true, value)
}

func Test_All_False(t *testing.T) {
	slice := []bool{true, false, true}
	value := iter.All(iter.IterSlice(&slice))
	assert.Equal(t, false, value)
}

func Test_Any_True(t *testing.T) {
	slice := []bool{true, false, true}
	value := iter.Any(iter.IterSlice(&slice))
	assert.Equal(t, true, value)
}

func Test_Any_False(t *testing.T) {
	slice := []bool{false, false, false}
	value := iter.Any(iter.IterSlice(&slice))
	assert.Equal(t, false, value)
}

func Test_Take(t *testing.T) {
	expectedValue := []int{1, 2}
	it := iter.IterSlice(&[]int{1, 2, 3})
	value := iter.Take[[]int](it, 2)

	expectedRem := []int{3}
	rem := iter.CollectToSlice[[]int](it)

	assert.Equal(t, expectedValue, value)
	assert.Equal(t, expectedRem, rem)
}
