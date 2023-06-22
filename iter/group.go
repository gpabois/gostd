package iter

type Groups[G comparable, T any, S ~[]T] map[G]S

func Group[S ~[]T, G comparable, T any](iter Iterator[T], key func(el T) G) Groups[G, T, S] {
	groups := make(Groups[G, T, S])

	for cursor := iter.Next(); cursor.IsSome(); cursor = iter.Next() {
		g := key(cursor.Expect())
		if _, ok := groups[g]; !ok {
			groups[g] = make([]T, 0)
		}
		groups[g] = append(groups[g], cursor.Expect())
	}

	return groups
}

func (groups Groups[G, T, S]) Iter() Iterator[KV[G, S]] {
	return IterMap(&groups)
}
