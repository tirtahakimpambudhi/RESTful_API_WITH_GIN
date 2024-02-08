package helper

import "go_gin/internal/exception"

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func CatchInternalServer() error {
	r := recover()
	if r != nil {
		err := r.(error)
		return exception.NewError(err, exception.ErrorInternalServer)
	}
	return nil
}
func NewCustomError(message, types error) error {
	if message != nil {
		return exception.NewError(message, types)
	}
	return nil
}
