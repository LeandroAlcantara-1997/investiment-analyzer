package response

import (
	"errors"
	"net/http"

	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/exception"
)

var restErrorMap = map[error]int{
	exception.ErrInternalServer: http.StatusInternalServerError,
	exception.ErrEntityNotFound: http.StatusNotFound,
	exception.ErrInvalidFile:    http.StatusBadRequest,
	exception.ErrInvalidRequest: http.StatusBadRequest,
}

func RestError(key error) (int, error) {
	var (
		err    *exception.ErrorWithTrace
		detail string
	)

	if errors.As(key, &err) {
		key = err.GetError()
		detail = err.Error()
	}
	code := restErrorMap[key]
	if code != 0 {
		return code, New(key, detail)
	}

	return http.StatusInternalServerError, New(exception.ErrInternalServer, err.Error())
}

type Exception struct {
	Key    string `json:"key"`
	Detail string `json:"detail"`
}

func New(err error, detail string) *Exception {
	return &Exception{
		Key:    err.Error(),
		Detail: detail,
	}
}

func (e *Exception) Error() string {
	return e.Key
}
