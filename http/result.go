package httputil

type HttpResult[T any] struct {
	Value T      `serde:"value"`
	Err   string `serde:"error"`
}

func (result HttpResult[T]) Success(value T) HttpResult[T] {
	res := HttpResult[T]{}
	res.Value = value
	res.Err = ""
	return res
}

func (result HttpResult[T]) Failed(err error) HttpResult[T] {
	res := HttpResult[T]{}
	res.Err = err.Error()
	return res
}
