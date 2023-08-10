package arena

import "github.com/gpabois/gostd/option"

type IHolder[T any] interface {
	Release()
	Copy() IHolder[T]
	Ref() *T
}

type IArena[ID any, T any] interface {
	New() ID
	Delete(ID)
	At(index ID) option.Option[IHolder[T]]
}

func With[ID any, T any](arena IArena[ID, T], id ID, fn func(ptr *T)) {
	arena.At(id).Then(func(h IHolder[T]) {
		defer h.Release()
		fn(h.Ref())
	})
}
