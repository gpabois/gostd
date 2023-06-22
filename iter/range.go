package iter

import (
	"github.com/gpabois/gostd/ops"
	"github.com/gpabois/gostd/option"
)

type rangeIterator[T ops.Addition] struct {
	r       ops.Interval[T]
	step    T
	current T
}

func (it *rangeIterator[T]) Next() option.Option[T] {
	if !ops.Within(it.current, it.r.Min(), it.r.Max()) {
		return option.None[T]()
	}

	res := option.Some(it.current)
	it.current = ops.Add(it.current, it.step)
	return res
}

func Range[T ops.Addition](min T, max T, step T) Iterator[T] {
	r := ops.NewInterval(min, max)
	return &rangeIterator[T]{
		r:       r,
		step:    step,
		current: r.Min(),
	}
}
