// +build go1.13

package errorx

import "reflect"

// todo godoc
// NB: Call to errors.As() converts any type of errorx error to any other type, therefore such calls may break semantics.
// Note than calls to errors.Is() do not suffer from the same issue.
// todo add errorx.As() ?
// todo make another type for passing into errors.As()?
// todo when fixed, add tests for wrap/decorate etc.
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

// todo godoc
func (e *Error) Is(target error) bool {
	typedTarget := Cast(target)
	if typedTarget == nil {
		return false
	}

	return e.IsOfType(typedTarget.Type())
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