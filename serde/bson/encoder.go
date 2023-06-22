package bson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gpabois/gostd/collection"
	"github.com/gpabois/gostd/result"
)

const (
	ENCODER_ROOT_STATE = iota
	ENCODER_ARRAY_STATE
	ENCODER_ARRAY_VALUE
	ENCODER_DOCUMENT_STATE
	ENCODER_DOCUMENT_KEY
	ENCODER_DOCUMENT_VALUE
)

type encoderState struct {
	typ      int
	buf      bytes.Buffer
	lastType byte
	key      string
	counter  int
}

type Encoder struct {
	states collection.Stack[encoderState]
	writer io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	enc := &Encoder{writer: w}
	enc.states.Push(encoderState{typ: ENCODER_ROOT_STATE})
	return enc
}

func (enc *Encoder) setType(b byte) {
	enc.states.Last().Expect().lastType = b
}

func (enc *Encoder) setKey(b []byte) {
	enc.states.Last().Expect().key = string(b)
}

func (enc *Encoder) Write(p []byte) (n int, err error) {
	return enc.states.Last().Expect().buf.Write(p)
}

func (enc *Encoder) EncodeInt64(value int64) result.Result[bool] {
	enc.setType(0x12)
	binary.Write(enc, binary.BigEndian, value)
	return result.Success(true)
}

func (enc *Encoder) EncodeFloat64(value float64) result.Result[bool] {
	enc.setType(0x01)
	binary.Write(enc, binary.BigEndian, value)
	return result.Success(true)
}

func (enc *Encoder) EncodeBool(value bool) result.Result[bool] {
	enc.setType(0x08)
	binary.Write(enc, binary.BigEndian, value)
	return result.Success(true)
}

func (enc *Encoder) EncodeString(value string) result.Result[bool] {
	enc.setType(0x02)
	binary.Write(enc, binary.BigEndian, value)
	return result.Success(true)
}

func (enc *Encoder) PushMap() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_DOCUMENT_STATE})
	return result.Success(true)
}

func (enc *Encoder) PushMapKey() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_DOCUMENT_KEY})
	return result.Success(true)
}

func (enc *Encoder) PushMapValue() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_DOCUMENT_VALUE})
	return result.Success(true)
}

func (enc *Encoder) PushArray() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_ARRAY_STATE})
	return result.Success(true)
}

func (enc *Encoder) PushArrayValue() result.Result[bool] {
	enc.states.Push(encoderState{typ: ENCODER_ARRAY_VALUE})
	return result.Success(true)
}

func (enc *Encoder) Pop() result.Result[bool] {
	last := enc.states.Pop().Expect()

	switch last.typ {
	// Flush the buffer into the writer
	case ENCODER_ROOT_STATE:
		last.buf.WriteTo(enc.writer)
	// Encode an array or a document
	case ENCODER_ARRAY_STATE, ENCODER_DOCUMENT_STATE:
		switch last.typ {
		case ENCODER_DOCUMENT_STATE:
			enc.setType(0x03)
		case ENCODER_ARRAY_STATE:
			enc.setType(0x04)
		}
		binary.Write(enc, binary.LittleEndian, int32(last.buf.Len()))
		last.buf.WriteTo(enc)
		enc.Write([]byte{0})
	// Encode a document key
	case ENCODER_DOCUMENT_KEY:
		enc.setKey(last.buf.Bytes())
	// Encode a document value
	case ENCODER_DOCUMENT_VALUE:
		binary.Write(enc, binary.LittleEndian, last.lastType) // Encode the value type
		enc.Write([]byte(enc.states.Last().Expect().key))     // Encode the key
		// A special care is required for string values
		switch last.lastType {
		case 0x02, 0x0C, 0x0D, 0x0E:
			binary.Write(enc, binary.LittleEndian, int32(last.buf.Len()))
		}
		last.buf.WriteTo(enc) // Fill the buffer with the value's data
	// Encode an array value
	case ENCODER_ARRAY_VALUE:
		// Generate an index
		index := fmt.Sprintf("%d", enc.states.Last().Expect().counter)
		enc.states.Last().Expect().counter++

		binary.Write(enc, binary.LittleEndian, last.lastType)
		enc.Write([]byte(index))
		// A special care is required for string values
		switch last.lastType {
		case 0x02, 0x0C, 0x0D, 0x0E:
			binary.Write(enc, binary.LittleEndian, int32(last.buf.Len()))
		}
		last.buf.WriteTo(enc)
	}

	return result.Success(true)
}
