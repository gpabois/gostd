package httputil

import (
	"context"
	"net/http"

	"github.com/gpabois/gostd/result"
	"github.com/gpabois/gostd/serde"
)

type HttpMiddleware interface {
	Handle(ctx context.Context, r *http.Request) result.Result[bool]
}

type DirectHandler[Request any, Response any] struct {
	endpoint    func(ctx context.Context, request Request) result.Result[Response]
	middlewares []HttpMiddleware
}

func (handler DirectHandler[Request, Response]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the content type of the request
	contentType := r.Header.Get("Content-Type")

	// Create the context
	ctx := context.Background()

	// Decode the request
	reqRes := serde.DeserializeFromReader[Request](r.Body, contentType)
	if reqRes.HasFailed() {
		EncodeError[Response](reqRes.UnwrapError(), w)
		return
	}

	// Pass it down to the endpoint
	respRes := handler.endpoint(ctx, reqRes.Expect())

	// Encode the result, we use the accept header
	accept := r.Header.Get("Accept")
	if accept == "" {
		accept = "application/json"
	}

	EncodeResult(w, respRes, contentType)
}

// A handler that allows a h(request) Result[Response]
func NewDirectHandler[Request any, Response any](endpoint func(ctx context.Context, request Request) result.Result[Response], middlewares ...HttpMiddleware) http.Handler {
	return DirectHandler[Request, Response]{
		endpoint,
		middlewares,
	}
}
