package errors

import (
	"fmt"
)

type AppError struct {
	Name string
	Err  error // wrap 元
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	err := e.Err
	if err == nil {
		return fmt.Sprintf("[%s]: unknown error", e.Name)
	}
	errMsg := err.Error()
	if len(errMsg) >= 1 && errMsg[0] == '[' {
		return fmt.Sprintf("[%s]%s", e.Name, errMsg)
	}

	return fmt.Sprintf("[%s]: %s", e.Name, errMsg)
}

func NewAppError(name string, err error) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Name: name,
		Err:  err,
	}
}

func (e *AppError) Unwrap() error { return e.Err }
