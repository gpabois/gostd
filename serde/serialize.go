package serde

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde/bson"
	"github.com/gpabois/gostd/serde/encoder"
	"github.com/gpabois/gostd/serde/json"
)

func getEncoderFromWriter(w io.Writer, contentType string) result.Result[encoder.Encoder] {
	switch contentType {
	case "application/bson":
		return result.Success[encoder.Encoder](bson.NewEncoder(w))
	case "application/json":
		return result.Success[encoder.Encoder](json.NewEncoder(w))
	default:
		return result.Failed[encoder.Encoder](errors.New(fmt.Sprintf("cannot unmarshal content-type '%s'", contentType)))
	}
}

// Serialize any value
func Serialize[T any](value T, contentType string) result.Result[[]byte] {
	var buf bytes.Buffer
	res := getEncoderFromWriter(&buf, contentType)
	if res.HasFailed() {
		return result.Result[[]byte]{}.Failed(res.UnwrapError())
	}
	encoder.Encode[T](res.Expect(), value)
	return result.Success(buf.Bytes())
}
