package cmd

import (
	"github.com/bitrise-io/bitrise-plugins-io/services"
)

// InputError ...
type InputError struct {
	Err string
}

func (e *InputError) Error() string {
	return e.Err
}

// NewInputError ...
func NewInputError(err string) error {
	return &InputError{
		Err: err,
	}
}

// RequestFailedError ...
type RequestFailedError struct {
	Response services.Response
}

func (e *RequestFailedError) Error() string {
	return e.Response.Error
}

// NewRequestFailedError ...
func NewRequestFailedError(response services.Response) error {
	return &RequestFailedError{
		Response: response,
	}
}
