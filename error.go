package healthy

import "errors"

type fatalError struct {
	err error
}

func Fatal(err error) error {
	return &fatalError{err: err}
}

func IsFatal(err error) bool {
	var f *fatalError
	return errors.As(err, &f)
}

func (e *fatalError) Error() string {
	return e.err.Error()
}

func (e *fatalError) Unwrap() []error {
	return []error{e.err}
}
