package validator

import "fmt"

// ValidationError represents a business rule validation failure
type ValidationError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error
func NewValidationError(code int, msg string, err error) *ValidationError {
	return &ValidationError{
		StatusCode: code,
		Message:    msg,
		Err:        err,
	}
}
