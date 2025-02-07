package exception

import (
	"errors"
)

var (
	ErrInternalServer = errors.New("unexpected error")
	ErrEntityNotFound = errors.New("entity not found")
	ErrInvalidFile    = errors.New("invalid format file")
	ErrInvalidRequest = errors.New("invalid request body")
)

type ErrorWithTrace struct {
	trace string
	err   error
}

func (e *ErrorWithTrace) GetError() error {
	return e.err
}

func (e *ErrorWithTrace) Error() string {
	return e.trace
}
func New(trace string, err error) *ErrorWithTrace {
	var errTrace *ErrorWithTrace
	if errors.As(err, &errTrace) {
		return &ErrorWithTrace{
			trace: errTrace.trace + trace,
			err:   err,
		}
	}

	return &ErrorWithTrace{
		trace: trace,
		err:   err,
	}
}
