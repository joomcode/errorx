// +build go1.13

package errorx

// Is returns true if and only if target is errorx error that passes errorx type check against current error.
// This behaviour is exactly the same as that of IsOfType().
// See also: errors.Is()
func (e *Error) Is(target error) bool {
	typedTarget := Cast(target)
	return typedTarget != nil && IsOfType(e, typedTarget.Type())
}

// From errors package: if e.Unwrap() returns a non-nil error w, then we say that e wraps w.
// Unwrap returns cause of current error in case it is wrapped transparently, nil otherwise.
// See also: errors.Unwrap()
func (e *Error) Unwrap() error {
	if e.cause != nil && e.transparent {
		return e.cause
	} else {
		return nil
	}
}