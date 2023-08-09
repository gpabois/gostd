package rb

import (
	"github.com/gpabois/gostd/arena"
	"github.com/gpabois/gostd/cmp"
	"github.com/gpabois/gostd/option"
)

// Red-black tree
type Tree[T any] struct {
	root  option.Option[id]
	nodes arena.IArena[node[T]]
}

func (tree *Tree[T]) Insert(value T) {
	nodeId := tree.nodes.New()
	node := tree.at(nodeId).Expect()
	node.value = value
	node.color = Red

	leafOpt := tree.getLeafId(value)
	node.parent = leafOpt

	// No leaf, so this is a new root
	if leafOpt.IsNone() {
		tree.root = option.Some(nodeId)
		return
	}

	tree.repairAfterInsertion(nodeId)
}

func (tree *Tree[T]) Remove(value T) {
	nodeIdOpt := tree.searchNodeByValue(value)
	if nodeIdOpt.IsNone() {
		return
	}

	nodeId := nodeIdOpt.Expect()
	node := tree.at(nodeId).Expect()

	parentOpt := tree.parent(nodeId)

	// Root leaf
	if node.isLeaf() && parentOpt.IsNone() {
		tree.root = option.None[id]()
		return
	}

	parent := parentOpt.Expect()
	if node.isLeaf() {
		parent.removeChild(nodeId)
		tree.nodes.Delete(nodeId)
	}
}

// Return the first node with the same value
func (tree *Tree[T]) searchNodeByValue(value T) option.Option[id] {
	currentId := tree.root
	for currentId.IsSome() {
		current := tree.at(currentId.Expect()).Expect()

		if current.isLeaf() {
			break
		}

		order := cmp.Cmp(value, current.value).Expect()

		if order == cmp.Equal {
			return currentId
		}

		currentId = current.getChildFromOrder(order)
	}
	return option.None[id]()
}

func (tree *Tree[T]) getLeafId(value T) option.Option[id] {
	currentId := tree.root
	for currentId.IsSome() {
		current := tree.at(currentId.Expect()).Expect()

		if current.isLeaf() {
			break
		}

		currentId = current.getChildFromOrder(cmp.Cmp(value, current.value).Expect())
	}
	return currentId
}

// Repair the tree
func (tree *Tree[T]) repairAfterInsertion(nodeId id) {
	// Case #1
	if tree.parentId(nodeId).IsNone() {
		tree.at(nodeId).Expect().color = Black

		return
	}

	// Case #2
	parentId := tree.parentId(nodeId).Expect()
	parent := tree.at(parentId).Expect()
	if parent.color == Black {
		return
	}

	// Case #3
	uncle := tree.uncle(nodeId)
	grandParentId := tree.grandParentId(nodeId).Expect()
	grandParent := tree.at(grandParentId).Expect()

	if uncle.IsSome() && uncle.Expect().color == Red {
		parent.color = Black
		uncle.Expect().color = Black

		grandParent.color = Red
		tree.repairAfterInsertion(grandParentId)

		return
	}

	// Case #4
	if leftRightChildId := tree.leftRightChildId(nodeId); leftRightChildId.IsSome() && leftRightChildId.Expect() == nodeId {
		tree.rotate(parentId, Left)
		nodeId = tree.at(nodeId).Expect().left.Expect()
	} else if rightLeftChildId := tree.rightLeftChildId(nodeId); rightLeftChildId.IsSome() && rightLeftChildId.Expect() == nodeId {
		tree.rotate(parentId, Right)
		nodeId = tree.at(nodeId).Expect().right.Expect()
	}

	// Case #5
	grandParentId = tree.grandParentId(nodeId).Expect()
	grandParent = tree.at(grandParentId).Expect()
	parent = tree.parent(nodeId).Expect()

	if parent.left.IsSome() && nodeId == parent.left.Expect() {
		tree.rotate(grandParentId, Right)
	} else {
		tree.rotate(grandParentId, Left)
	}

	parent.color = Black
	grandParent.color = Red
}

// Rotate a node, either left or right
func (tree *Tree[T]) rotate(xNodeId id, dir direction) {
	yNodeId := tree.childId(xNodeId, 1-dir)
	yNode := tree.at(yNodeId.Expect()).Expect()
	xNode := tree.at(xNodeId).Expect()

	pivot := tree.childId(yNodeId.Expect(), dir)

	xNode.setChildId(1-dir, pivot)
	yNode.setChildId(dir, option.Some(xNodeId))

	yNode.parent = xNode.parent
	if xNode.parent.IsNone() {
		tree.root = yNodeId
	}
}

// Get the node behind the ID
func (tree *Tree[T]) at(nodeId id) option.Option[*node[T]] {
	return tree.nodes.At(nodeId)
}

func (tree *Tree[T]) childId(nodeId id, dir direction) option.Option[id] {
	return option.Flatten(option.Map(tree.at(nodeId), func(n *node[T]) option.Option[id] {
		return n.getChildId(dir)
	}))
}

// Get the parent id, if any.
func (tree *Tree[T]) parentId(nodeId id) option.Option[id] {
	return option.Flatten(option.Map(tree.at(nodeId), func(n *node[T]) option.Option[id] {
		return n.parent
	}))
}

// Get the parent.
func (tree *Tree[T]) parent(nodeId id) option.Option[*node[T]] {
	return option.Flatten(option.Map(tree.parentId(nodeId), func(nodeId id) option.Option[*node[T]] {
		return tree.at(nodeId)
	}))
}

// n->right->left
func (tree *Tree[T]) rightLeftChildId(nodeId id) option.Option[id] {
	return option.Flatten(option.Map(tree.at(nodeId), func(c *node[T]) option.Option[id] {
		right := c.right.Expect()
		return option.Flatten(option.Map(tree.at(right), func(c *node[T]) option.Option[id] {
			return c.left
		}))
	}))
}

// n->left->right
func (tree *Tree[T]) leftRightChildId(nodeId id) option.Option[id] {
	return option.Flatten(option.Map(tree.at(nodeId), func(c *node[T]) option.Option[id] {
		left := c.left.Expect()
		return option.Flatten(option.Map(tree.at(left), func(c *node[T]) option.Option[id] {
			return c.right
		}))
	}))
}

// Get the uncle id.
func (tree *Tree[T]) uncleId(nodeId id) option.Option[id] {
	return option.Flatten(option.Map(tree.parent(nodeId), func(parent *node[T]) option.Option[id] {
		return parent.getOtherChildId(nodeId)
	}))
}

// Get the uncle.
func (tree *Tree[T]) uncle(nodeId id) option.Option[*node[T]] {
	return option.Flatten(option.Map(tree.uncleId(nodeId), func(uncleId id) option.Option[*node[T]] {
		return tree.at(uncleId)
	}))
}

// Get the grand parent id, if any.
func (tree *Tree[T]) grandParentId(nodeId id) option.Option[id] {
	return option.Flatten(option.Map(tree.parent(nodeId), func(node *node[T]) option.Option[id] {
		return node.parent
	}))
}

// Get the grand parent.
func (tree *Tree[T]) grandParent(nodeId id) option.Option[*node[T]] {
	return option.Flatten(option.Map(tree.grandParentId(nodeId), func(grandParentId id) option.Option[*node[T]] {
		return tree.at(grandParentId)
	}))
}
