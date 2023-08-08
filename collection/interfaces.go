package collection

import "github.com/gpabois/gostd/option"

type Sizable interface {
	Length() uint
}

type Indexable[T any] interface {
	At(index uint) option.Option[T]
	RefAt(index uint) option.Option[*T]
}
