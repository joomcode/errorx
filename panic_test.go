package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPanic(t *testing.T) {

	defer func() {
		r := recover()
		require.NotNil(t, r)
		output := fmt.Sprintf("%v", r)

		require.Contains(t, output, "errorx.funcWithErr()", output)
	}()

	Panic(funcWithErr())
}

func TestPanicErrorx(t *testing.T) {

	defer func() {
		r := recover()
		require.NotNil(t, r)
		output := fmt.Sprintf("%v", r)

		require.Contains(t, output, "awful", output)
		// original error was non-errorx, no use of adding panic callstack to message
		require.NotContains(t, output, "errorx.funcWithBadPanic()", output)
	}()

	funcWithBadPanic()
}

func TestPanicRecover(t *testing.T) {

	defer func() {
		r := recover()
		require.NotNil(t, r)

		err, ok := ErrorFromPanic(r)
		require.True(t, ok)

		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "errorx.funcWithErr()", output)
		require.Contains(t, output, "bad", output)
		require.True(t, IsOfType(err, testType))
	}()

	Panic(funcWithErr())
}

func TestPanicRecoverNoTrace(t *testing.T) {

	defer func() {
		r := recover()
		require.NotNil(t, r)

		err, ok := ErrorFromPanic(r)
		require.True(t, ok)

		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "errorx.funcWithErrNoTrace()", output)
		require.Contains(t, output, "errorx.funcWithPanicNoTrace()", output)
		require.Contains(t, output, "silent", output)
		require.True(t, IsOfType(err, testType))
	}()

	funcWithPanicNoTrace()
}

func TestPanicRecoverNoErrorx(t *testing.T) {

	defer func() {
		r := recover()
		require.NotNil(t, r)

		err, ok := ErrorFromPanic(r)
		require.True(t, ok)

		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "errorx.funcWithBadErr()", output)
		require.Contains(t, output, "errorx.funcWithBadPanic()", output)
		require.Contains(t, output, "awful", output)
		require.False(t, IsOfType(err, testType))
	}()

	funcWithBadPanic()
}

func funcWithErr() error {
	return testType.New("bad")
}

func funcWithPanicNoTrace() {
	Panic(funcWithErrNoTrace())
}

func funcWithErrNoTrace() error {
	return testTypeSilent.New("silent")
}

func funcWithBadPanic() {
	Panic(funcWithBadErr())
}

func funcWithBadErr() error {
	return errors.New("awful")
}
