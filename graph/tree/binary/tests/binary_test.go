package tests

import (
	"testing"

	"github.com/gpabois/gostd/graph/tree/binary"
	"github.com/gpabois/gostd/option"
	"github.com/stretchr/testify/assert"
)

func Test_NewNode(t *testing.T) {
	tree := binary.Default(100)

	nodeId := tree.NewNode()
	assert.True(t, tree.NodeExists(nodeId))
}

/*
IN:

	     ┌─┐
	   ┌─┤Q├─┐
	   │ └─┘ │
	  ┌┴┐   ┌┴┐
	┌─┤P├─┐ │C│
	│ └─┘ │ └─┘

┌┴┐   ┌┴┐
│A│   │B│
└─┘   └─┘

# Perform a rotation of Q to the right

OUT :

	  ┌─┐
	┌─┤P├─┐
	│ └─┘ │

┌┴┐   ┌┴┐
│A│  ┌┤Q├─┐
└─┘  │└─┘ │

	┌┴┐  ┌┴┐
	│B│  │C│
	└─┘  └─┘
*/
func Test_RotateRight(t *testing.T) {
	tree := binary.Default(5)

	qId := tree.NewNode()
	pId := tree.NewNode()
	aId := tree.NewNode()
	bId := tree.NewNode()
	cId := tree.NewNode()

	// Set Q children, left : P, right: C
	tree.WithNode(qId, func(n *binary.Node) {
		n.SetChild(binary.Left, option.Some(pId))
		n.SetChild(binary.Right, option.Some(cId))
	})

	// Set P children, left : A, right : B
	tree.WithNode(pId, func(n *binary.Node) {
		n.SetChild(binary.Left, option.Some(aId))
		n.SetChild(binary.Right, option.Some(bId))
	})

	// Perform the node rotation.
	binary.Rotate[binary.Node, *binary.Node](&tree, qId, binary.Right)

	// P should be the new root
	assert.Equal(t, tree.GetRoot(), option.Some(pId))

	// P should have A as left child, and Q as the right child
	assert.True(t, tree.NodeExists(pId))
	tree.WithNode(pId, func(p *binary.Node) {
		assert.Equal(t, p.GetChild(binary.Left), option.Some(aId))
		assert.Equal(t, p.GetChild(binary.Right), option.Some(qId))
	})

	// Q should have B as the left child, and C as the right child
	assert.True(t, tree.NodeExists(qId))
	tree.WithNode(qId, func(q *binary.Node) {
		assert.Equal(t, q.GetChild(binary.Left), option.Some(bId))
		assert.Equal(t, q.GetChild(binary.Right), option.Some(cId))
	})
}

/*
IN:

	  ┌─┐
	┌─┤P├─┐
	│ └─┘ │

┌┴┐   ┌┴┐
│A│  ┌┤Q├─┐
└─┘  │└─┘ │

	┌┴┐  ┌┴┐
	│B│  │C│
	└─┘  └─┘

# Perform a rotation of P to the left

OUT :

	     ┌─┐
	   ┌─┤Q├─┐
	   │ └─┘ │
	  ┌┴┐   ┌┴┐
	┌─┤P├─┐ │C│
	│ └─┘ │ └─┘

┌┴┐   ┌┴┐
│A│   │B│
└─┘   └─┘
*/
func Test_RotateLeft(t *testing.T) {
	tree := binary.Default(5)

	qId := tree.NewNode()
	pId := tree.NewNode()
	aId := tree.NewNode()
	bId := tree.NewNode()
	cId := tree.NewNode()

	// Set P children, left : A, right: Q
	tree.WithNode(pId, func(n *binary.Node) {
		n.SetChild(binary.Left, option.Some(aId))
		n.SetChild(binary.Right, option.Some(qId))
	})

	// Set Q children, left : B, right : C
	tree.WithNode(qId, func(n *binary.Node) {
		n.SetChild(binary.Left, option.Some(bId))
		n.SetChild(binary.Right, option.Some(cId))
	})

	// Perform the node rotation.
	binary.Rotate[binary.Node, *binary.Node](&tree, pId, binary.Left)

	// Q should be the new root
	assert.Equal(t, tree.GetRoot(), option.Some(qId))

	// Q should have P as left child, and C as the right child
	assert.True(t, tree.NodeExists(qId))
	tree.WithNode(qId, func(q *binary.Node) {
		assert.Equal(t, q.GetChild(binary.Left), option.Some(pId))
		assert.Equal(t, q.GetChild(binary.Right), option.Some(cId))
	})

	// P should have A as the left child, and B as the right child
	assert.True(t, tree.NodeExists(pId))
	tree.WithNode(pId, func(p *binary.Node) {
		assert.Equal(t, p.GetChild(binary.Left), option.Some(aId))
		assert.Equal(t, p.GetChild(binary.Right), option.Some(bId))
	})
}
