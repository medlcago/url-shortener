package http

import (
	"net/http"
)

type Exception struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *Exception) Error() string {
	return e.Message
}

func (e *Exception) SetMessage(message string) *Exception {
	e.Message = message
	return e
}

func (e *Exception) SetStatus(status int) *Exception {
	e.Status = status
	return e
}

func NewException(status int, message string) *Exception {
	return &Exception{
		Status:  status,
		Message: message,
	}
}

func WithMessage(message string) *Exception {
	return NewException(http.StatusInternalServerError, message)
}

func WithStatus(status int) *Exception {
	return NewException(status, http.StatusText(status))
}
