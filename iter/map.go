package iter

import (
	"reflect"

	"github.com/gpabois/gostd/option"
)

type MapIterator[K comparable, V any, M ~map[K]V] struct {
	keyIt Iterator[K]
	inner *M
}

func (it *MapIterator[K, V, M]) Next() option.Option[KV[K, V]] {
	key := it.keyIt.Next()
	if key.IsNone() {
		return option.None[KV[K, V]]()
	}

	if val, ok := (*it.inner)[key.Expect()]; ok {
		return option.Some(KV[K, V]{key: key.Expect(), value: val})
	}

	return option.None[KV[K, V]]()
}

func IterMap[K comparable, V any, M ~map[K]V](m *M) Iterator[KV[K, V]] {
	keys := reflect.ValueOf(*m).MapKeys()
	return &MapIterator[K, V, M]{
		keyIt: Map(IterSlice(keys), func(value reflect.Value) K { return value.Interface().(K) }),
		inner: m,
	}
}
