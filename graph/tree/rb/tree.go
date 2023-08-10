package rb

import (
	"github.com/gpabois/gostd/arena"
	"github.com/gpabois/gostd/cfg"
	"github.com/gpabois/gostd/cmp"
	"github.com/gpabois/gostd/option"
)

// Red-black tree
// A balanced binary search tree, with insertion, search and remove operations at a ln(N) time.
type Tree[T any] struct {
	root  option.Option[id]
	nodes arena.IArena[id, node[T]]
}

type TreeOptions[T any] struct {
	nodes option.Option[arena.IArena[id, node[T]]]
}

func NewTree[T any](options ...cfg.Configurator[TreeOptions[T]]) Tree[T] {
	opt := TreeOptions[T]{}
	cfg.Apply(&opt, options)
	return Tree[T]{
		nodes: opt.nodes.UnwrapOr(func() arena.IArena[id, node[T]] {
			a := arena.NewSimpleArena[node[T]]()
			return &a
		}),
	}
}

// Search the value in the tree, or return the upper bound.
func (tree *Tree[T]) SearchEqualOrUpperBound(value T) option.Option[T] {
	currentId := tree.root
	res := option.None[T]()
	found := false

	for currentId.IsSome() && !found {
		tree.withNode(currentId.Expect(), func(n *node[T]) {
			order := cmp.Cmp(value, n.value).Expect()

			// Found you
			if order == cmp.Equal {
				res = option.Some[T](n.value)
				found = true
			} else { // Keep the max value
				res.Then(func(max T) {
					res = option.Some(cmp.Max(max, n.value).Expect())
				}).Else(func() {
					res = option.Some(n.value)
				})
			}

			currentId = n.getChildFromOrder(order)
		})
	}
	return res
}

// Search the value in the tree.
func (tree *Tree[T]) Search(value T) option.Option[T] {
	currentId := tree.root
	res := option.None[T]()

	for currentId.IsSome() && res.IsNone() {
		tree.withNode(currentId.Expect(), func(n *node[T]) {
			order := cmp.Cmp(value, n.value).Expect()
			// Found you
			if order == cmp.Equal {
				res = option.Some[T](n.value)
			}
			currentId = n.getChildFromOrder(order)
		})
	}
	return res
}

// Insert a value in the tree.
func (tree *Tree[T]) Insert(value T) {
	// Create a new node
	nodeId := tree.nodes.New()

	// Initialise the node
	tree.withNode(nodeId, func(node *node[T]) {
		node.id = nodeId
		node.value = value
		node.color = Red

		tree.searchLeafId(value).Then(func(leafId id) {
			node.parent = option.Some(leafId)
		}).Else(func() {
			tree.root = option.Some(nodeId)
		})
	})

	// Repair the tree afterwards
	tree.repairAfterInsertion(nodeId)
}

func (tree *Tree[T]) searchLeafId(value T) option.Option[id] {
	currentId := tree.root
	found := false
	for currentId.IsSome() && !found {
		tree.withNode(currentId.Expect(), func(n *node[T]) {
			if n.isLeaf() {
				found = true
				return
			}

			currentId = n.getChildFromOrder(cmp.Cmp(value, n.value).Expect())
		})
	}
	return currentId
}

// Repair the tree after insertion
func (tree *Tree[T]) repairAfterInsertion(nodeId id) {
	tree.withNode(nodeId, func(n *node[T]) {
		// Case #1
		// The inserted node has not parent,
		// thus the node should be black.
		if n.parent.IsNone() {
			n.color = Black
			return
		}

		tree.withParent(nodeId, func(parent *node[T]) {
			// Case #2
			// If the parent's color is black, the tree is valid.
			// The black-height does not change as the new node is red.
			if parent.color == Black {
				return
			}

			// Parent is red, thus the tree is invalid.
			// What to do depends of the uncle's color.
			tree.withGrandParent(nodeId, func(grandParent *node[T]) {
				tree.withUncle(nodeId, func(uncle *node[T]) {
					// Case #3,
					// Uncle is red.
					// So the parent and the uncle must be black, and the grand-parent red.
					// We repair the tree from the grand-parent perspective.
					if uncle.color == Red {
						parent.color = Black
						uncle.color = Black
						grandParent.color = Red
						tree.repairAfterInsertion(grandParent.id)
						return
					}

					// Case #4
					// Uncle is black
					// We need to perform rotation, the parent becomes the child of the node
					if tree.isLeftRightChild(grandParent.id, nodeId) {
						tree.rotate(parent.id, Left)
						nodeId = n.left.Expect()
					} else if tree.isRightLeftChild(grandParent.id, nodeId) {
						tree.rotate(parent.id, Right)
						nodeId = n.right.Expect()
					}

					// Case #5
					// Parent becomes the grand-parent, and the grand-parent is the new uncle.
					// Parent becomes black, and the new grand-parent becomes red
					tree.withParent(nodeId, func(parent *node[T]) {
						tree.withGrandParent(nodeId, func(grandParent *node[T]) {
							if tree.isLeftChild(parent.id, nodeId) {
								tree.rotate(grandParent.id, Right)
							} else {
								tree.rotate(grandParent.id, Left)
							}

							parent.color = Black
							grandParent.color = Red
						})
					})

				})
			})
		})
	})
}

func (tree *Tree[T]) isRightChild(parentId id, nodeId id) bool {
	result := false
	tree.withNode(parentId, func(parent *node[T]) {
		parent.right.Then(func(rightId id) {
			result = rightId == nodeId
		})
	})
	return result
}

func (tree *Tree[T]) isLeftChild(parentId id, nodeId id) bool {
	result := false
	tree.withNode(parentId, func(parent *node[T]) {
		parent.left.Then(func(leftId id) {
			result = leftId == nodeId
		})
	})
	return result
}

func (tree *Tree[T]) isRightLeftChild(grandParentId id, nodeId id) bool {
	result := false
	tree.withNode(grandParentId, func(grandParent *node[T]) {
		grandParent.right.Then(func(rightId id) {
			tree.withNode(rightId, func(right *node[T]) {
				right.left.Then(func(leftId id) {
					result = leftId == nodeId
				})
			})
		})
	})
	return result
}

func (tree *Tree[T]) isLeftRightChild(grandParentId id, nodeId id) bool {
	result := false
	tree.withNode(grandParentId, func(grandParent *node[T]) {
		grandParent.left.Then(func(leftId id) {
			tree.withNode(leftId, func(left *node[T]) {
				left.right.Then(func(rightId id) {
					result = rightId == nodeId
				})
			})
		})
	})
	return result
}

// Rotate a node, either left or right
func (tree *Tree[T]) rotate(xNodeId id, dir direction) {
	tree.withNode(xNodeId, func(xNode *node[T]) {
		xNode.getChildId(1 - dir).Then(func(yNodeId id) {
			tree.withNode(yNodeId, func(yNode *node[T]) {
				pivot := yNode.getChildId(dir)

				xNode.setChildId(1-dir, pivot)
				yNode.setChildId(dir, option.Some(xNodeId))

				yNode.parent = xNode.parent
				if xNode.parent.IsNone() {
					tree.root = option.Some(yNodeId)
				}
			})
		})
	})
}

func (tree *Tree[T]) withNode(nodeId id, fn func(*node[T])) {
	arena.With(tree.nodes, nodeId, fn)
}

func (tree *Tree[T]) withParent(nodeId id, fn func(*node[T])) {
	parentIdOpt := option.None[id]()

	tree.withNode(nodeId, func(node *node[T]) {
		parentIdOpt = node.parent
	})

	parentIdOpt.Then(func(parentId id) {
		tree.withNode(parentId, fn)
	})
}

func (tree *Tree[T]) withGrandParent(nodeId id, fn func(*node[T])) {
	grandParentIdOpt := option.None[id]()

	tree.withParent(nodeId, func(node *node[T]) {
		grandParentIdOpt = node.parent
	})

	grandParentIdOpt.Then(func(grandParentId id) {
		tree.withNode(grandParentId, fn)
	})
}

func (tree *Tree[T]) withUncle(nodeId id, fn func(*node[T])) {
	tree.withParent(nodeId, func(parent *node[T]) {
		parent.getOtherChildId(nodeId).Then(func(uncleId id) {
			tree.withNode(uncleId, fn)
		})
	})
}
