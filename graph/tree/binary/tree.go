package binary

import (
	"github.com/gpabois/gostd/mmu/zoo"
	"github.com/gpabois/gostd/option"
)

type Step = int8
type Direction = int8

type NodePtr[N any] interface {
	*N
	INode
}

type ITree[N any, PN NodePtr[N]] interface {
	SetRoot(option.Option[NodeId])
	GetRoot() option.Option[NodeId]

	TryNewNode() option.Option[NodeId]
	NewNode() NodeId
	NodeExists(nodeId NodeId) bool

	WithNode(nodeId NodeId, fn func(n PN))
}

// A binary tree
type Tree[N any, PN NodePtr[N]] struct {
	root  option.Option[NodeId]
	nodes zoo.Zoo[NodeId, N, PN]
}

// Create a new binary tree
func Default(capacity uint) Tree[Node, *Node] {
	z := zoo.Default[Node](capacity)
	return Tree[Node, *Node]{
		nodes: &z,
	}
}

// Tries to create a new node
func (tree *Tree[N, PN]) TryNewNode() option.Option[NodeId] {
	return tree.nodes.TryNew()
}

// Creates a new node, panics if cannot.
func (tree *Tree[N, PN]) NewNode() NodeId {
	return tree.nodes.New()
}

func (tree *Tree[N, PN]) NodeExists(nodeId NodeId) bool {
	return tree.nodes.Exists(nodeId)
}

// Get a reference to the node
func (tree *Tree[N, PN]) WithNode(id NodeId, fn func(n PN)) {
	tree.nodes.With(id, fn)
}

func (tree *Tree[N, PN]) GetRoot() option.Option[NodeId] {
	return tree.root
}

func (tree *Tree[N, PN]) SetRoot(root option.Option[NodeId]) {
	tree.root = root
}

// Perform a node rotation of a binary tree (Source: https://en.wikipedia.org/wiki/Tree_rotation)
func Rotate[N any, PN NodePtr[N]](tree ITree[N, PN], xNodeId NodeId, dir Direction) {
	tree.WithNode(xNodeId, func(xNode PN) {
		xNode.GetChild(1 - dir).Then(func(yNodeId NodeId) {
			tree.WithNode(yNodeId, func(yNode PN) {
				pivot := yNode.GetChild(dir)
				xNode.SetChild(1-dir, pivot)
				yNode.SetChild(dir, option.Some(xNodeId))

				yNode.SetParent(xNode.GetParent())
				xNode.GetParent().Else(func() {
					tree.SetRoot(option.Some(yNodeId))
				})
			})
		})
	})
}

func FollowStep[N any, PN NodePtr[N]](tree ITree[N, PN], nodeId NodeId, path *[]Step) option.Option[NodeId] {
	if len(*path) == 0 {
		return option.Some(nodeId)
	}
	step := (*path)[0]
	*path = (*path)[1:]
	res := option.None[NodeId]()
	tree.WithNode(nodeId, func(n PN) {
		res = n.GetChild(step)
	})
	return res
}

func Follow[N any, PN NodePtr[N]](tree ITree[N, PN], start NodeId, path []Step) option.Option[NodeId] {
	cursorId := option.Some(start)

	for len(path) > 0 {
		// Path is not empty, and we failed to get a next step...
		if cursorId.IsNone() {
			return cursorId
		}
		cursorId = FollowStep(tree, cursorId.Expect(), &path)
	}

	return cursorId
}

// Check is the node behind the path is the same as the node passed as the argument
func Is[N any, PN NodePtr[N]](tree ITree[N, PN], start NodeId, nodeId NodeId, path []Step) bool {
	targetOpt := Follow(tree, start, path)

	return option.Map(targetOpt, func(targetId NodeId) bool {
		return targetId == nodeId
	}).UnwrapOrZero()
}
