package iter_tests

import (
	"testing"

	"github.com/gpabois/gostd/iter"
	"github.com/stretchr/testify/assert"
)

func Test_ChunkEvery(t *testing.T) {
	c1 := iter.CollectToSlice[iter.Chunk[int]](iter.Range(0, 9, 1))
	c2 := iter.CollectToSlice[iter.Chunk[int]](iter.Range(10, 19, 1))
	c3 := iter.CollectToSlice[iter.Chunk[int]](iter.Range(20, 29, 1))
	c4 := iter.CollectToSlice[iter.Chunk[int]](iter.Range(30, 30, 1))

	expectedChunks := []iter.Chunk[int]{c1, c2, c3, c4}

	chunks := iter.CollectToSlice[[]iter.Chunk[int]](iter.ChunkEvery(iter.Range(0, 30, 1), 10))

	assert.Equal(t, expectedChunks, chunks)
}
