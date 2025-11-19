package status

import (
	"errors"
	"fmt"
)

var (
	ErrStatusNotFound = errors.New("status not found")
)

type ServiceError struct {
    StatusCode int
    Message    string
    Err        error
}

func (e *ServiceError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *ServiceError) Unwrap() error {
    return e.Err
}

func NewServiceError(statusCode int, message string, err error) *ServiceError {
    return &ServiceError{
        StatusCode: statusCode,
        Message:    message,
        Err:        err,
    }
}