// This file implements the helper interfaces for the errors package introduced in Go 1.13
// +build go1.13

package errorx

import "errors"

// As implements the helper interface for errors.As.
// If v is a pointer to a pointer to Error, *v will be set to this error.
// Otherwise, it traverses the error chain until it reaches a non-*Error error and calls errors.As on that error and the
// target.
func (e *Error) As(v interface{}) bool {
	if target, ok := v.(**Error); ok {
		*target = e
		return true
	}

	cause := e
	for cause != nil {
		next := cause.Cause()
		cause = Cast(next)
		if cause == nil && next != nil {
			return errors.As(next, v)
		}
	}
	return false
}

// Is implements the helper interface for errors.Is.
// If target is an *Error, we check each non-transparent Error in the error chain and compare the types using IsOfType.
// If target is not an *Error, we traverse the error chain until we reach a non-Error error and use errors.Is to check
// if that error Is the same as the target error.
func (e *Error) Is(target error) bool {
	var t *Type
	if cast := Cast(target); cast != nil {
		t = cast.Type()
	}

	cause := e
	for cause != nil {
		if t != nil && !cause.transparent && cause.errorType.IsOfType(t) {
			return true
		}
		next := cause.Cause()
		cause = Cast(next)
		// We check if t == nil as well because a non-Error error will fail the errors.Is check anyways if target is
		// an Error.
		if t == nil && cause == nil && next != nil {
			return errors.Is(next, target)
		}
	}
	return false
}

// Unwrap implements the helper interface for errors.Unwrap.
// It will return the first, non-transparent error in the error chain, or the final error if it is not an *Error.
// It is not an inversion of the Wrap method.
func (e *Error) Unwrap() error {
	cause := e
	for cause != nil {
		next := cause.Cause()
		cause = Cast(next)
		if cause == nil {
			return next
		} else if !cause.transparent {
			return cause
		}
	}
	return nil
}
