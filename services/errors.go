package services

// RequestFailedError ...
type RequestFailedError struct {
	Response Response
}

func (e *RequestFailedError) Error() string {
	return e.Response.Error
}

// NewRequestFailedError ...
func NewRequestFailedError(response Response) error {
	return &RequestFailedError{
		Response: response,
	}
}
