package rb

import (
	"github.com/gpabois/gostd/cmp"
	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/option"
)

const (
	Red = byte(iota)
	Black
)

const (
	Left = byte(iota)
	Right
)

type color = byte
type direction = byte
type id = uint

type node[T any] struct {
	id    id
	color color
	value T

	parent option.Option[id]
	left   option.Option[id]
	right  option.Option[id]
}

func (n *node[T]) Iter() iter.Iterator[id] {
	c := []id{}
	n.left.Then(func(cid id) {
		c = append(c, cid)
	})
	n.right.Then(func(cid id) {
		c = append(c, cid)
	})

	return iter.IterSlice(&c)
}

func (n *node[T]) removeChild(nodeId id) {
	nodeIdOpt := option.Some(nodeId)

	if nodeIdOpt == n.left {
		n.left = option.None[id]()
	} else if nodeIdOpt == n.right {
		n.right = option.None[id]()
	}
}

func (n *node[T]) getChildFromOrder(order cmp.Order) option.Option[id] {
	switch order {
	case cmp.Less:
		return n.left
	default:
		return n.right
	}
}

func (n *node[T]) isLeaf() bool {
	return n.left.IsNone() && n.right.IsNone()
}

func (n *node[T]) getOtherChildId(nodeId id) option.Option[id] {
	nodeIdOpt := option.Some(nodeId)

	if nodeIdOpt == n.left {
		return n.right
	} else {
		return n.left
	}
}

func (n *node[T]) getChildId(dir direction) option.Option[id] {
	switch dir {
	case Left:
		return n.left
	case Right:
		return n.right
	default:
		return option.None[id]()
	}
}

func (n *node[T]) setChildId(dir direction, nodeId option.Option[id]) {
	switch dir {
	case Left:
		n.left = nodeId
	case Right:
		n.right = nodeId
	}
}
