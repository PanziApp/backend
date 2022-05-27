package domain

import "fmt"

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid error: %s", e.Err.Error())
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

type InternalError struct {
	Err error
	// TODO: Add origin field here and to all internal errors.
}

func (e InternalError) Error() string {
	return fmt.Sprintf("internal error")
}

func (e InternalError) Unwrap() error {
	return e.Err
}

type ServiceError struct {
	Name string
	Err  error
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("service error")
}

func (e ServiceError) Unwrap() error {
	return e.Err
}
