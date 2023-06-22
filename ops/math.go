package ops

import "golang.org/x/exp/constraints"

type Number interface {
	int | float32 | float64
}

type Addition interface {
	~int | ~float32 | ~float64
}

type Interval[T constraints.Ordered] struct {
	begin T
	end   T
}

func NewInterval[T constraints.Ordered](min T, max T) Interval[T] {
	return Interval[T]{
		begin: min,
		end:   max,
	}
}

func (itv Interval[T]) Within(val T) bool {
	return val >= itv.Min() && val <= itv.Max()
}

func (itv Interval[T]) Max() T {
	return Max(itv.begin, itv.end)
}

func (itv Interval[T]) Min() T {
	return Min(itv.begin, itv.end)
}

func Add[T Addition](values ...T) T {
	var acc T
	for _, value := range values {
		acc = acc + value
	}
	return acc
}

func Max[T constraints.Ordered](acc T, values ...T) T {
	for _, value := range values {
		if value > acc {
			acc = value
		}
	}
	return acc
}

func Max2[T constraints.Ordered](a, b T) T {
	return Max(a, b)
}

func Min[T constraints.Ordered](acc T, values ...T) T {
	for _, value := range values {
		if value < acc {
			acc = value
		}
	}
	return acc
}

func Within[T constraints.Ordered](val T, min T, max T) bool {
	return NewInterval(min, max).Within(val)
}

func Bounds[T constraints.Ordered](val T, min T, max T) T {
	return Min(Max(val, min), max)
}

func Min2[T constraints.Ordered](a, b T) T {
	return Min(a, b)
}

func Add2[T Addition](a, b T) T {
	return Add(a, b)
}

func IsTrue(b bool) bool {
	return b
}

func IsFalse(b bool) bool {
	return !b
}
