package graph

import (
	"github.com/gpabois/gostd/iter"
	"github.com/gpabois/gostd/option"
)

type depthFirstIterator[N TreeBranch[N]] struct {
	nodes iter.Iterator[*N]
}

func (d *depthFirstIterator[N]) Next() option.Option[*N] {
	nodeOption := d.nodes.Next()
	if nodeOption.IsNone() {
		return option.None[*N]()
	}
	node := nodeOption.Expect()
	d.nodes = iter.Chain(d.nodes, (*node).IterChildren())
	return nodeOption
}
