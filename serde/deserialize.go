package serde

import (
	"bytes"
	"io"
	"reflect"

	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde/decoder"
	"github.com/gpabois/gostd/serde/json"
)

func getDecoderFromReader(r io.Reader, contentType string) result.Result[decoder.Decoder] {
	switch contentType {
	case "application/json":
		return result.Success[decoder.Decoder](json.NewDecoder(r))
	default:
		return result.Failed[decoder.Decoder](NewUnhandledContentType(contentType))
	}
}

func getDecoderFromBytes(b []byte, contentType string) result.Result[decoder.Decoder] {
	buf := bytes.NewBuffer(b)
	return getDecoderFromReader(buf, contentType)
}

func DeserializeFromReaderInto[T any](r io.Reader, contentType string, ptr *T) result.Result[bool] {
	res := getDecoderFromReader(r, contentType)
	if res.HasFailed() {
		return result.Result[bool]{}.Failed(res.UnwrapError())
	}
	return decoder.DecodeInto(res.Expect(), ptr)
}

func DeserializeFromReader[T any](r io.Reader, contentType string) result.Result[T] {
	res := getDecoderFromReader(r, contentType)
	if res.HasFailed() {
		return result.Result[T]{}.Failed(res.UnwrapError())
	}
	return decoder.Decode[T](res.Expect())
}

func DeserializeInto[T any](value []byte, contentType string, ptr *T) result.Result[bool] {
	res := getDecoderFromBytes(value, contentType)
	if res.HasFailed() {
		return result.Result[bool]{}.Failed(NewDeserializeError(res.UnwrapError()))
	}
	return decoder.DecodeInto(res.Expect(), ptr)
}

func Deserialize[T any](value []byte, contentType string) result.Result[T] {
	res := getDecoderFromBytes(value, contentType)
	if res.HasFailed() {
		return result.Result[T]{}.Failed(NewDeserializeError(res.UnwrapError()))
	}
	return decoder.Decode[T](res.Expect())
}

func Reflect_Deserialize[T any](value []byte, contentType string, typ reflect.Type) result.Result[any] {
	res := getDecoderFromBytes(value, contentType)
	if res.HasFailed() {
		return result.Result[any]{}.Failed(NewDeserializeError(res.UnwrapError()))
	}
	return decoder.Reflect_Decode(res.Expect(), typ)
}
