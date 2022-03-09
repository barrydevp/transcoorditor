package util

import (
	"fmt"
)

func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

func Errorw(dst error, src error) error {
	return fmt.Errorf("%w: %v", src, dst)
}
