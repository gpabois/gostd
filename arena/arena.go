package arena

import "github.com/gpabois/gostd/option"

type IArena[T any] interface {
	New() uint
	Delete(uint)
	At(index uint) option.Option[*T]
}
