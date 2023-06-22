package ops_tests

import (
	"testing"

	"github.com/gpabois/gostd/ops"
	"github.com/stretchr/testify/assert"
)

func Test_Interval(t *testing.T) {
	itv := ops.NewInterval(10, 0)
	assert.Equal(t, itv.Min(), 0)
	assert.Equal(t, itv.Max(), 10)
}

func Test_Within(t *testing.T) {
	assert.True(t, ops.Within(1, 1, 2))
	assert.False(t, ops.Within(0, 1, 2))
	assert.False(t, ops.Within(3, 1, 2))
}

func Test_Bounds(t *testing.T) {
	assert.Equal(t, 2, ops.Bounds(3, 1, 2))
	assert.Equal(t, 1, ops.Bounds(0, 1, 2))
}

func Test_IsTrue_IsFalse(t *testing.T) {
	assert.Equal(t, true, ops.IsTrue(true))
	assert.Equal(t, false, ops.IsTrue(false))
	assert.Equal(t, true, ops.IsFalse(false))
	assert.Equal(t, false, ops.IsFalse(true))
}

func Test_Max(t *testing.T) {
	assert.Equal(t, 3, ops.Max(3, 1, 2))
}

func Test_Min(t *testing.T) {
	assert.Equal(t, 1, ops.Min(3, 1, 2))
}

func Test_Add(t *testing.T) {
	assert.Equal(t, 6, ops.Add(1, 2, 3))
}
