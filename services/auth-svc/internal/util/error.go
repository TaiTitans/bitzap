package util

import (
	"fmt"
)

// NewError message
func NewError(message string) error {
	return fmt.Errorf("%s", message)
}

// NewErrorf format
func NewErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// WrapError wrap error message
func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf wrap error format
func WrapErrorf(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}
