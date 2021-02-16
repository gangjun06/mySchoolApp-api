package errors

type (
	ErrorCode string
	Error     struct {
		Code    ErrorCode
		Message string
	}
)

const (
	ErrBadRequest ErrorCode = "BAD_REQUEST"
)

var (
	ErrorCodeDescription = map[ErrorCode]string{
		ErrBadRequest: "This is an incorrect request",
	}
)

func New(code ErrorCode, message string) error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Description() string {
	return ErrorCodeDescription[e.Code]
}
