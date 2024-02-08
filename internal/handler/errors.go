package handler

import (
	"errors"
	"go_gin/internal/exception"
	"net/http"
)

type ResponseErrors struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewResponseErrors(err error) *ResponseErrors {
	var responseErrors ResponseErrors // Initialize without using a pointer
	typeErrors := &exception.Error{}

	if errors.As(err, &typeErrors) {
		responseErrors.Message = typeErrors.MessageError().Error()
		typeErr := typeErrors.TypeError()
		switch typeErr {
		case exception.ErrorNotFound:
			responseErrors.Status = http.StatusNotFound
		case exception.ErrorInternalServer:
			responseErrors.Status = http.StatusInternalServerError
		case exception.ErrorBadRequest:
			responseErrors.Status = http.StatusBadRequest
		case exception.ErrorConflict:
			responseErrors.Status = http.StatusConflict
		case exception.ErrorForbidden:
			responseErrors.Status = http.StatusForbidden
		case exception.ErrorUnauthorized:
			responseErrors.Status = http.StatusUnauthorized
		}
	}
	return &responseErrors // Return a pointer to ResponseErrors
}
