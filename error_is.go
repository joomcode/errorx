// +build go1.13

package errorx

import "errors"

// Is implements the helper interface for errors.Is
// It first checks if target is an *Error and uses IsOfType to check if e Is target.
// If target is not an *Error, it instead iterates through the error chain until it reaches a non-*Error error (or nil)
// and calls errors.Is on that error and target.
func (e *Error) Is(target error) bool {
	// Check if target is an *Error, if so we can use IsOfType.
	t := Cast(target)
	if t != nil {
		return e.IsOfType(t.Type())
	}

	// target is not an *Error, so iterate through the error chain until nil or a non-*Error error is reached and use
	// errors.Is to check against that.
	cause := e
	for cause != nil {
		next := Cast(cause.Cause())
		if next == nil && cause.Cause() != nil {
			return errors.Is(cause.Cause(), target)
		}
		cause = next
	}
	return false
}