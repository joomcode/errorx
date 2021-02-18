package errorx

import "fmt"

// todo godoc
type typeCheckTarget struct {
	err *Error
	assignabilityBlocker interface{}
}

func (t typeCheckTarget) Error() string {
	return t.err.Error()
}

func (t typeCheckTarget) Format(s fmt.State, verb rune) {
	t.err.Format(s, verb)
}

func (t *typeCheckTarget) AsError() *Error {
	return t.err
}

var _ error = typeCheckTarget{}
