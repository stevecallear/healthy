package healthy

import "errors"

type fatalError struct {
	err error
}

// Fatal wraps the supplied error to indicate that it is fatal.
// When a check returns a fatal error retry execution will be aborted.
func Fatal(err error) error {
	return &fatalError{err: err}
}

// IsFatal returns true if the supplied error is fatal.
func IsFatal(err error) bool {
	var f *fatalError
	return errors.As(err, &f)
}

// Error returns the inner error message.
func (e *fatalError) Error() string {
	return e.err.Error()
}

// Unwrap returns the inner error.
func (e *fatalError) Unwrap() []error {
	return []error{e.err}
}
