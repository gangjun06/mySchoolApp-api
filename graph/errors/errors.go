package errors

import (
	"strings"
)

type (
	ErrorCode string
	Error     struct {
		Code    ErrorCode
		Message string
	}
)

const (
	ErrBadRequest    ErrorCode = "BAD_REQUEST"
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrPasswordWrong ErrorCode = "PASSWORD_WRONG"
	ErrAuth          ErrorCode = "AUTH"
	ErrTooManyReq    ErrorCode = "TOO_MANY_REQ"
	ErrServer        ErrorCode = "ERR_SERVER"
)

var (
	ErrorCodeDescription = map[ErrorCode]string{
		ErrBadRequest:    "This is an incorrect request",
		ErrNotFound:      "can not find item",
		ErrPasswordWrong: "password is incorrect",
		ErrAuth:          "error while auth",
		ErrTooManyReq:    "server get too many requests. try again later",
		ErrServer:        "error server",
	}
)

func New(code ErrorCode, message string) error {
	return &Error{Code: code, Message: message}
}

func Parse(err error) (error, bool) {
	spl := strings.Split(err.Error(), ": ")
	if len(spl) != 2 || ErrorCodeDescription[ErrorCode(spl[0])] == "" {
		return nil, false
	}
	return New(ErrorCode(spl[0]), spl[1]), true
}

func (e *Error) Error() string {
	return string(e.Code) + ": " + e.Message
}

func (e *Error) Description() string {
	return ErrorCodeDescription[e.Code]
}
