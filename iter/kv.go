package iter

// Key-Value
type KV[K any, V any] struct {
	key   K
	value V
}

func (pair KV[T, U]) Key() T {
	return pair.key
}

func (pair KV[T, U]) Value() U {
	return pair.value
}

func (pair KV[T, U]) GetFirst() T {
	return pair.key
}

func (pair KV[T, U]) GetSecond() U {
	return pair.value
}
