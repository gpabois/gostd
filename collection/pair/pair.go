package collection

type Pair[T, U any] struct {
	First  T
	Second U
}

func (pair Pair[T, U]) GetFirst() T {
	return pair.First
}

func (pair Pair[T, U]) GetSecond() U {
	return pair.Second
}

type IPair[T, U any] interface {
	GetFirst() T
	GetSecond() U
}
