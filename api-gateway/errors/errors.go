package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest          = fmt.Errorf("invalid request")
	ErrUnauthorized        = fmt.Errorf("request unauthorized")
	ErrNotFound            = fmt.Errorf("resource not found")
	ErrDuplicatedKey       = fmt.Errorf("request violate unique constraint")
	ErrInternalServerError = fmt.Errorf("internal server error")
)

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func GetStatusCode(err error) (code int) {
	switch {
	case errors.Is(err, ErrBadRequest):
		code = http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		code = http.StatusUnauthorized
	case errors.Is(err, ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, ErrDuplicatedKey):
		code = http.StatusConflict
	default:
		code = http.StatusInternalServerError
	}

	return
}
