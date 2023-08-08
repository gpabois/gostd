package graph

import "github.com/gpabois/gostd/iter"

type TreeBranch[N any] interface {
	IterChildren() iter.Iterator[*N]
}
