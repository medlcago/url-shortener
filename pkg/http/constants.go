package http

const (
	ErrBadRequest          = "bad request"
	ErrEmailAlreadyExists  = "user with given email already exists"
	ErrInvalidCredentials  = "invalid credentials"
	ErrInternalServer      = "internal server error"
	ErrInvalidToken        = "invalid token"
	ErrTokenAlreadyExpired = "token already expired"
	ErrNotFound            = "not found"
)

var (
	BadRequest          = NewException(400, ErrBadRequest)
	InvalidCredentials  = NewException(401, ErrInvalidCredentials)
	InternalServerError = NewException(500, ErrInternalServer)
	ExistsEmailError    = NewException(409, ErrEmailAlreadyExists)
	InvalidToken        = NewException(401, ErrInvalidToken)
	TokenExpired        = NewException(401, ErrTokenAlreadyExpired)
	NotFound            = NewException(404, ErrNotFound)
)
