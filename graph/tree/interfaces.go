package tree

import "github.com/gpabois/gostd/option"

type SearchTree[T any] interface {
	Search(value T) option.Option[T]
	Insert(value T)
	Remove(value T)
}
