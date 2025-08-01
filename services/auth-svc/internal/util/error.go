package util

import (
	"fmt"
)

// NewError creates a new error with message
func NewError(message string) error {
	return fmt.Errorf("%s", message)
}

// NewErrorf creates a new formatted error
func NewErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// WrapError wraps an error with additional message
func WrapError(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf wraps an error with formatted message
func WrapErrorf(err error, format string, args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}
