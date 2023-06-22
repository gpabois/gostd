package httputil

import (
	"context"
	"net/http"

	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde"
)

func EncodeResult[T any](w http.ResponseWriter, res result.Result[T], contentType string) {

}

func EncodeSuccess[T any](ctx context.Context, w http.ResponseWriter, response T) error {
	contentType := "application/json"
	encodeResult := serde.Serialize(
		HttpResult[T]{}.Success(response),
		contentType,
	)

	if encodeResult.IsSuccess() {
		w.Write(encodeResult.Expect())
	}

	return encodeResult.UnwrapError()
}

func EncodeError[T any](err error, w http.ResponseWriter) {
	contentType := "application/json"
	w.Header().Set("Content-Type", contentType)
	encodeResult := serde.Serialize(
		HttpResult[T]{}.Failed(err),
		contentType,
	)

	if encodeResult.IsSuccess() {
		w.Write(encodeResult.Expect())
	}
}
