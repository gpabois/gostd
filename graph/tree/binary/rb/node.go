package rb

import (
	"github.com/gpabois/gostd/graph/tree/binary"
)

const (
	Red = byte(iota)
	Black
)

const (
	Left = byte(iota)
	Right
)

type color = byte
type direction = byte
type id = uint

// A Red-black Node
type Node[T any] struct {
	color color
	value T

	binary.Node
}
