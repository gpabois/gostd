package tests

import (
	"fmt"
	"testing"

	"github.com/gpabois/gostd/graph/tree/rb"
	"github.com/stretchr/testify/assert"
)

func Test_Insertion(t *testing.T) {
	tree := rb.NewTree[int]()

	fmt.Println(tree)

	tree.Insert(10)
	assert.True(t, tree.Search(10).IsSome(), fmt.Sprintf("Should have found 10"))

	tree.Insert(5)
	assert.True(t, tree.Search(5).IsSome(), fmt.Sprintf("Should have found 5"))

	tree.Insert(11)
	assert.True(t, tree.Search(11).IsSome(), fmt.Sprintf("Should have found 11"))
}
