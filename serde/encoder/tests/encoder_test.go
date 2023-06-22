package encoder_tests

import (
	"testing"

	"github.com/gpabois/gostd/option"
	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde/encoder"
	"github.com/gpabois/gostd/serde/encoder/mocks"
)

type simpleStruct struct {
	Opt option.Option[bool]
	Val int `serde:"val"`
}

func Test_EncodeInt64(t *testing.T) {
	enc := mocks.NewEncoder(t)

	enc.EXPECT().EncodeInt64(int64(10)).Return(result.Success(true))

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, 10)
}

func Test_EncodeBool(t *testing.T) {
	enc := mocks.NewEncoder(t)

	enc.EXPECT().EncodeBool(true).Return(result.Success(true))

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, true)
}

func Test_EncodeFloat64(t *testing.T) {
	enc := mocks.NewEncoder(t)

	enc.EXPECT().EncodeFloat64(1.02).Return(result.Success(true))

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, 1.02)
}

func Test_EncodeString(t *testing.T) {
	enc := mocks.NewEncoder(t)

	enc.EXPECT().EncodeString("hello").Return(result.Success(true))

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, "hello")
}

func Test_EncodeSlice(t *testing.T) {
	enc := mocks.NewEncoder(t)

	enc.EXPECT().PushArray().Return(result.Success(true))
	enc.EXPECT().PushArrayValue().Return(result.Success(true))

	enc.EXPECT().EncodeString("hello").Return(result.Success(true))

	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Array Value
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Array

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, []string{"hello"})
}

func Test_EncodeMap(t *testing.T) {
	enc := mocks.NewEncoder(t)

	enc.EXPECT().PushMap().Return(result.Success(true))

	enc.EXPECT().PushMapKey().Return(result.Success(true))
	enc.EXPECT().EncodeString("hello").Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Key

	enc.EXPECT().PushMapValue().Return(result.Success(true))
	enc.EXPECT().EncodeString("world").Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Value

	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Value

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, map[string]string{"hello": "world"})
}

func Test_EncodeStruct_OptionNone(t *testing.T) {
	v := simpleStruct{Val: 10}

	enc := mocks.NewEncoder(t)

	enc.EXPECT().PushMap().Return(result.Success(true))

	enc.EXPECT().PushMapKey().Return(result.Success(true))
	enc.EXPECT().EncodeString("val").Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Key

	enc.EXPECT().PushMapValue().Return(result.Success(true))
	enc.EXPECT().EncodeInt64(int64(10)).Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Value

	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Value

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, v)
}

func Test_EncodeStruct_OptionSome(t *testing.T) {
	v := simpleStruct{Val: 10, Opt: option.Some(true)}

	enc := mocks.NewEncoder(t)

	enc.EXPECT().PushMap().Return(result.Success(true))

	enc.EXPECT().PushMapKey().Return(result.Success(true))
	enc.EXPECT().EncodeString("Opt").Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Key

	enc.EXPECT().PushMapValue().Return(result.Success(true))
	enc.EXPECT().EncodeBool(true).Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Value

	enc.EXPECT().PushMapKey().Return(result.Success(true))
	enc.EXPECT().EncodeString("val").Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Key

	enc.EXPECT().PushMapValue().Return(result.Success(true))
	enc.EXPECT().EncodeInt64(int64(10)).Return(result.Success(true))
	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Value

	enc.EXPECT().Pop().Return(result.Success(true)) // Pop Map Value

	// Pop the root state of the encoder at the end of the process.
	enc.EXPECT().Pop().Return(result.Success(true))

	encoder.Encode(enc, v)
}
