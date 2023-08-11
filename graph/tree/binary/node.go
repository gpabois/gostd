package binary

import (
	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/option"
)

type NodeId = uint

const Left = 0
const Right = 1

type INode interface {
	SetParent(option.Option[NodeId])
	GetParent() option.Option[NodeId]

	GetChild(dir Direction) option.Option[NodeId]
	SetChild(dir Direction, nodeId option.Option[NodeId])

	iter.Iterable[NodeId]
}

type Node struct {
	parent   option.Option[NodeId]
	children [2]option.Option[NodeId]
}

func (n *Node) SetParent(parentId option.Option[NodeId]) {
	n.parent = parentId
}

func (n *Node) GetParent() option.Option[NodeId] {
	return n.parent
}

func (n *Node) OtherChild(nodeId NodeId) option.Option[NodeId] {
	if n.children[0] == option.Some(nodeId) {
		return n.children[1]
	} else {
		return n.children[0]
	}
}

func (n *Node) IsLeaf() bool {
	return n.children[0].IsNone() && n.children[1].IsNone()
}

func (n *Node) GetChild(dir Direction) option.Option[NodeId] {
	return n.children[ops.Bounds(dir, 0, 1)]
}

func (n *Node) SetChild(dir Direction, nodeId option.Option[NodeId]) {
	n.children[ops.Bounds(dir, 0, 1)] = nodeId
}

func (n *Node) Iter() iter.Iterator[NodeId] {
	cOptIt := iter.IterSlice(n.children[:])

	return iter.Map(
		iter.Filter(cOptIt, option.IsSome[NodeId]),
		option.Expect[NodeId],
	)
}
