package errs

import "net/http"

type Error interface {
	StatusCode() int
	Status() string
	Message() string
}

type CustomError struct {
	ErrStatusCode int    `json:"status_code"`
	ErrStatus     string `json:"status"`
	ErrMessage    string `json:"message"`
	Data          any    `json:"data"`
}

func (c *CustomError) StatusCode() int {
	return c.ErrStatusCode
}

func (c *CustomError) Status() string {
	return c.ErrStatus
}

func (c *CustomError) Message() string {
	return c.ErrMessage
}

func NewInternalServerError(msg string) Error {
	return &CustomError{
		ErrStatusCode: http.StatusInternalServerError,
		ErrStatus:     http.StatusText(http.StatusInternalServerError),
		ErrMessage:    msg,
		Data:          nil,
	}
}

func NewNotFoundError(msg string) Error {
	return &CustomError{
		ErrStatusCode: http.StatusNotFound,
		ErrStatus:     http.StatusText(http.StatusNotFound),
		ErrMessage:    msg,
		Data:          nil,
	}
}

func NewBadRequestError(msg string) Error {
	return &CustomError{
		ErrStatusCode: http.StatusBadRequest,
		ErrStatus:     http.StatusText(http.StatusBadRequest),
		ErrMessage:    msg,
		Data:          nil,
	}
}

func NewConflictError(msg string) Error {
	return &CustomError{
		ErrStatusCode: http.StatusConflict,
		ErrStatus:     http.StatusText(http.StatusConflict),
		ErrMessage:    msg,
		Data:          nil,
	}
}

func NewUnprocessableEntityError(msg string) Error {
	return &CustomError{
		ErrStatusCode: http.StatusUnprocessableEntity,
		ErrStatus:     http.StatusText(http.StatusUnprocessableEntity),
		ErrMessage:    msg,
		Data:          nil,
	}
}
