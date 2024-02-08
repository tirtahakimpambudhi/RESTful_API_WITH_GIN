package exception

import "errors"

var (
	ErrorInternalServer = errors.New("INTERNAL SERVER ERROR")
	ErrorBadRequest     = errors.New("BAD REQUEST")
	ErrorConflict       = errors.New("CONFLICT : DUPLICATE KEY VIOLATION")
	ErrorNotFound       = errors.New("NOT FOUND")
	ErrorUnauthorized   = errors.New("UNAUTHORIZED")
	ErrorForbidden      = errors.New("FORBIDDEN")
)

type Error struct {
	messageError error
	typeError    error
}

func (e *Error) MessageError() error {
	return e.messageError
}

func (e *Error) TypeError() error {
	return e.typeError
}

func (e *Error) Error() string {
	return errors.Join(e.typeError, e.messageError).Error()
}

func NewError(messageError error, typeError error) *Error {
	return &Error{messageError: messageError, typeError: typeError}
}
