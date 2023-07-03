package norm

import (
	"github.com/gpabois/gostd/collection"
	"github.com/gpabois/gostd/result"
)

const (
	root = iota
	mapState
	mapKeyState
	mapValueState
	arrayState
	arrayValueState
)

type state struct {
	typ int
	r0  any
	r1  any
	r2  any
}

type Encoder struct {
	states collection.Stack[state]
	ret    any
}

func NewEncoder() *Encoder {
	e := Encoder{}
	e.states.Push(state{typ: root})
	return &e
}

func (e *Encoder) Encoded() any {
	return e.ret
}

func (e *Encoder) EncodeInt64(value int64) result.Result[bool] {
	e.current().r0 = value
	return result.Success(true)

}

func (e *Encoder) EncodeFloat64(value float64) result.Result[bool] {
	e.current().r0 = value
	return result.Success(true)

}

func (e *Encoder) EncodeBool(value bool) result.Result[bool] {
	e.current().r0 = value
	return result.Success(true)
}

func (e *Encoder) EncodeString(value string) result.Result[bool] {
	e.current().r0 = value
	return result.Success(true)
}

func (e *Encoder) PushArray() result.Result[bool] {
	e.states.Push(state{typ: arrayState, r0: []any{}})
	return result.Success(true)
}
func (e *Encoder) PushArrayValue() result.Result[bool] {
	e.states.Push(state{typ: arrayValueState})
	return result.Success(true)
}
func (e *Encoder) PushMap() result.Result[bool] {
	e.states.Push(state{typ: mapState, r0: map[string]any{}})
	return result.Success(true)
}
func (e *Encoder) PushMapKey() result.Result[bool] {
	e.states.Push(state{typ: mapKeyState})
	return result.Success(true)
}
func (e *Encoder) PushMapValue() result.Result[bool] {
	e.states.Push(state{typ: mapValueState})
	return result.Success(true)
}

func (e *Encoder) Pop() result.Result[bool] {
	last := e.states.Pop().Expect()

	switch last.typ {
	case mapKeyState:
		e.current().r1 = last.r0
	case mapValueState:
		e.current().r2 = last.r0
		e.reduceKeyValue()
	case arrayValueState:
		e.current().r1 = last.r0
		e.reduceValue()
	case root:
		e.ret = last.r0
	default:
		e.current().r0 = last.r0
	}

	return result.Success(true)
}

func (e *Encoder) current() *state {
	return e.states.Last().Expect()
}

func (e *Encoder) currentMap() map[string]any {
	return e.current().r0.(map[string]any)
}

func (e *Encoder) currentArray() []any {
	return e.current().r0.([]any)
}

func (e *Encoder) reduceKeyValue() {
	e.currentMap()[e.current().r1.(string)] = e.current().r2
}

func (e *Encoder) reduceValue() {
	e.current().r0 = append(e.currentArray(), e.current().r1)
}
