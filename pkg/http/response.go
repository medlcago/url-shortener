package http

type Response[T any] struct {
	Ok   bool   `json:"ok"`
	Err  string `json:"error,omitempty"`
	Data T      `json:"data,omitempty"`
}

func NewResponse[T any](ok bool, err string, data ...T) *Response[T] {
	res := &Response[T]{
		Ok:  ok,
		Err: err,
	}
	if len(data) > 0 {
		res.Data = data[0]
	}
	return res
}

func NewErrorResponse(err string) *Response[string] {
	return NewResponse[string](false, err)
}
