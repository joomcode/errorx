// +build go1.13

package errorx

import "errors"

func isOfType(err error, t *Type) bool {
	e := burrowForTyped(err)
	return e != nil && e.IsOfType(t)
}

func (e *Error) isOfType(t *Type) bool {
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
