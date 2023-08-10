package rb

import (
	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/option"
)

func (tree *Tree[T]) Iter() iter.Iterator[T] {
	it := newTreeIterator(tree)
	return &it
}

type treeIterator[T any] struct {
	tree  *Tree[T]
	stack iter.Iterator[id]
}

func newTreeIterator[T any](tree *Tree[T]) treeIterator[T] {
	c := []id{}
	tree.root.Then(func(nodeId id) {
		c = append(c, nodeId)
	})

	return treeIterator[T]{
		tree:  tree,
		stack: iter.IterSlice(&c),
	}
}

func (it *treeIterator[T]) Next() option.Option[T] {
	return option.Chain(it.stack.Next(), func(nodeId id) option.Option[T] {
		el := option.None[T]()
		it.tree.withNode(nodeId, func(node *node[T]) {
			it.stack = iter.Chain(it.stack, node.Iter())
			el = option.Some(node.value)
		})
		return el
	})
}
