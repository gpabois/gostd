package json

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/gpabois/gostd/collection"
	"github.com/gpabois/gostd/result"
)

type encoderState struct {
	typ     int
	counter int
	buffer  bytes.Buffer
}

const (
	ENCODER_ROOT_STATE = iota
	ENCODER_MAP_STATE
	ENCODER_MAP_KEY_STATE
	ENCODER_MAP_VALUE_STATE
	ENCODER_ARRAY_STATE
	ENCODER_ARRAY_VALUE_STATE
)

type Encoder struct {
	states collection.Stack[encoderState]
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	enc := &Encoder{writer: w}
	enc.states.Push(encoderState{typ: ENCODER_ROOT_STATE})
	return enc
}

func (enc *Encoder) EncodeInt64(value int64) result.Result[bool] {
	return enc.write(fmt.Sprintf("%d", value))
}
func (enc *Encoder) EncodeFloat64(value float64) result.Result[bool] {
	return enc.write(fmt.Sprintf("%f", value))
}
func (enc *Encoder) EncodeBool(value bool) result.Result[bool] {
	if value {
		return enc.write("true")
	} else {
		return enc.write("false")
	}
}
func (enc *Encoder) EncodeString(value string) result.Result[bool] {
	return enc.write(fmt.Sprintf("\"%s\"", value))
}

func (enc *Encoder) write(s string) result.Result[bool] {
	_, err := enc.writer.Write([]byte(s))

	if err != nil {
		return result.Failed[bool](err)
	}

	return result.Success(true)
}

func (enc *Encoder) PushArray() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_ARRAY_STATE})
	return enc.write("[")
}

func (enc *Encoder) PushArrayValue() result.Result[bool] {
	currentStateRes := enc.states.Last().IntoResult(errors.New("no state"))
	if currentStateRes.HasFailed() {
		return result.Result[bool]{}.Failed(currentStateRes.UnwrapError())
	}

	currentState := currentStateRes.Expect()

	if currentState.typ != ENCODER_ARRAY_STATE {
		return result.Result[bool]{}.Failed(errors.New("expecting the encoder to be in a map state"))
	}

	if currentState.counter > 0 {
		enc.write(",")
	}
	currentState.counter++
	enc.states.Push(encoderState{typ: ENCODER_ARRAY_VALUE_STATE})

	return result.Success(true)
}

func (enc *Encoder) PushMap() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_MAP_STATE})
	enc.write("{")
	return result.Success(true)
}

func (enc *Encoder) PushMapKey() result.Result[bool] {
	currentStateRes := enc.states.Last().IntoResult(errors.New("no state"))

	if currentStateRes.HasFailed() {
		return result.Result[bool]{}.Failed(currentStateRes.UnwrapError())
	}

	currentState := currentStateRes.Expect()

	if currentState.typ != ENCODER_MAP_STATE {
		return result.Result[bool]{}.Failed(errors.New("expecting the encoder to be in a map state"))
	}

	if currentState.counter > 0 {
		if res := enc.write(","); res.HasFailed() {
			return result.Result[bool]{}.Failed(res.UnwrapError())
		}
	}

	currentState.counter++
	enc.states.Push(encoderState{typ: ENCODER_MAP_KEY_STATE})
	return result.Success(true)
}

func (enc *Encoder) PushMapValue() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_MAP_VALUE_STATE})
	enc.write(":")
	return result.Success(true)
}

func (enc *Encoder) Pop() result.Result[bool] {
	state := enc.states.Pop().IntoResult(errors.New("no state was popped"))

	if state.HasFailed() {
		return result.Result[bool]{}.Failed(state.UnwrapError())
	}

	switch state.Expect().typ {
	case ENCODER_MAP_STATE:
		return enc.write("}")
	case ENCODER_ARRAY_STATE:
		return enc.write("]")
	case ENCODER_ARRAY_VALUE_STATE:
	case ENCODER_MAP_VALUE_STATE:
	case ENCODER_MAP_KEY_STATE:
		return result.Success(true)
	default:
		return result.Failed[bool](errors.New("invalid encoder state"))
	}

	return result.Success(true)
}
