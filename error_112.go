// +build !go1.13

package errorx

func isOfType(err error, t *Type) bool {
	e := Cast(err)
	return e != nil && e.IsOfType(t)
}

func (e *Error) isOfType(t *Type) bool {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.errorType.IsOfType(t)
		}

		cause = Cast(cause.Cause())
	}

	return false
}


