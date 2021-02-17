// +build go1.13

package errorx

// todo godoc
func (e *Error) As(target interface{}) bool {
	targetError, ok := target.(*error)
	if !ok {
		return false
	}

	if !e.Is(*targetError) {
		return false
	}

	// todo inject
	return true
}

// todo godoc
func (e *Error) Is(target error) bool {
	typedTarget := Cast(target)
	if typedTarget == nil {
		return false
	}

	return e.IsOfType(typedTarget.Type()) // todo test with parents
}

// todo godoc
// If e.Unwrap() returns a non-nil error w, then we say that e wraps w.
func (e *Error) Unwrap() error {
	if e.cause != nil && e.transparent {
		return e.cause
	} else {
		return nil
	}
}