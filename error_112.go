// +build !go1.13

package errorx

// IsOfType is a type check for errors.
// Returns true either if both are of exactly the same type, or if the same is true for one of current type's ancestors.
// For an error that does not have an errorx type, returns false.
func IsOfType(err error, t *Type) bool {
	e := Cast(err)
	return e != nil && e.IsOfType(t)
}

// IsOfType is a proper type check for an errorx-based errors.
// It takes the transparency and error types hierarchy into account,
// so that type check against any supertype of the original cause passes.
func (e *Error) IsOfType(t *Type) bool {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.errorType.IsOfType(t)
		}

		cause = Cast(cause.Cause())
	}

	return false
}


