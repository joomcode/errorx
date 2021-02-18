// +build go1.13

package errorx

import "reflect"

// todo add errorx.As()?
// todo make another 'go error type' for passing into errors.As()?
// todo when fixed, add tests for wrap/decorate etc.
// As checks if target is of the same type as current error and, if true, sets target to this error value.
// NB: Call to errors.As() converts any type of errorx error to any other type,
// therefore such calls are currently unsafe for errorx errors and will likely break semantics.
// Note than calls to errors.Is() do not suffer from the same issue.
func (e *Error) As(target interface{}) bool {
	targetError, ok := target.(*error)
	if !ok {
		return false
	}

	if !e.Is(*targetError) {
		return false
	}

	targetVal := reflect.ValueOf(target)
	targetVal.Elem().Set(reflect.ValueOf(e))
	return true
}

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