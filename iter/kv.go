package iter

// Key-Value
type KV[K any, V any] struct {
	Key   K
	Value V
}

func (pair KV[T, U]) GetFirst() T {
	return pair.Key
}

func (pair KV[T, U]) GetSecond() U {
	return pair.Value
}
