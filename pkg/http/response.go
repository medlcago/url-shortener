package http

type MetaData map[string]any

type Response[T any] struct {
	Ok   bool     `json:"ok"`
	Err  string   `json:"error,omitempty"`
	Data T        `json:"data"`
	Meta MetaData `json:"meta,omitempty"`
}

func NewResponse[T any](ok bool, err string, data T, meta ...MetaData) *Response[T] {
	var m MetaData
	if len(meta) > 0 {
		m = meta[0]
	}

	res := &Response[T]{
		Ok:   ok,
		Err:  err,
		Data: data,
		Meta: m,
	}
	return res
}

func Error(err string) *Response[any] {
	return NewResponse[any](false, err, nil)
}

func OK[T any](data T, meta ...MetaData) *Response[T] {
	return NewResponse[T](true, "", data, meta...)
}
