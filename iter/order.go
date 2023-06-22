package iter

import (
	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/option"
	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](it Iterator[T]) option.Option[T] {
	firstOpt := it.Next()

	if firstOpt.IsNone() {
		return option.None[T]()
	}

	first := firstOpt.Expect()

	return option.Some(Reduce(it, ops.Max2[T], first))
}

func Min[T constraints.Ordered](it Iterator[T]) option.Option[T] {
	firstOpt := it.Next()

	if firstOpt.IsNone() {
		return option.None[T]()
	}

	first := firstOpt.Expect()

	return option.Some(Reduce(it, ops.Min2[T], first))
}
