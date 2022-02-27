package util

import (
    "fmt"
)

func NewError(format string, a ...interface{}) error {
    return fmt.Errorf(format, a...)
}
