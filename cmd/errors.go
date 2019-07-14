package cmd

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
