// +build go1.13

package errorx

import "errors"

// IsOfType is a type check for errors.
// Returns true either if both are of exactly the same type, or if the same is true for one of current type's ancestors.
// For an error that does not have an errorx type, returns false unless it wraps another error of errorx type.
func IsOfType(err error, t *Type) bool {
	e := burrowForTyped(err)
	return e != nil && e.IsOfType(t)
}


// IsOfType is a proper type check for an errorx-based errors.
// It takes the transparency and error types hierarchy into account,
// so that type check against any supertype of the original cause passes.
// It also tolerates non-errorx errors in chain if those errors support go 1.13 errors unwrap.
func (e *Error) IsOfType(t *Type) bool {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.errorType.IsOfType(t)
		}

		cause = burrowForTyped(cause.Cause())
	}

	return false
}

// burrowForTyped returns either the first *Error in unwrap chain or nil
func burrowForTyped(err error) *Error {
	raw := err
	for raw != nil {
		typed := Cast(raw)
		if typed != nil {
			return typed
		}

		raw = errors.Unwrap(raw)
	}

	return nil
}


