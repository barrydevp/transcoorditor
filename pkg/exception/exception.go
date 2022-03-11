package exception

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound           = errors.New("NOT_FOUND")
	ErrInvalidArgument    = errors.New("INVALID_ARGUMENT")
	ErrPreconditionFailed = errors.New("PRECONDITION_FAILED")
	ErrAborted            = errors.New("ABORTED")
)


func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

func Errorw(dst error, src error) error {
	return fmt.Errorf("%w: %v", src, dst)
}
