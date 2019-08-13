package errorx

import (
	"errors"
	"fmt"
	"testing"
	"time"

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
		require.Contains(t, output, "errorx.funcWithBadPanic()", output)
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

func TestPanicChain(t *testing.T) {
	ch0 := make(chan error, 1)
	ch1 := make(chan error, 1)

	go doMischief(ch1)
	go doMoreMischief(ch0, ch1)

	select {
	case err := <-ch0:
		require.Error(t, err)
		require.False(t, IsOfType(err, AssertionFailed))
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "mischiefProper", output)
		require.Contains(t, output, "mischiefAsPanic", output)
		require.Contains(t, output, "doMischief", output)
		require.Contains(t, output, "handleMischief", output)
		require.NotContains(t, output, "doMoreMischief", output) // stack trace is only enhanced in Panic, not in user code
		t.Log(output)
	case <-time.After(time.Second):
		require.Fail(t, "expected error")
	}
}

func doMoreMischief(ch0 chan error, ch1 chan error) {
	defer func() {
		if e := recover(); e != nil {
			err, ok := ErrorFromPanic(e)
			if ok {
				ch0 <- Decorate(err, "hop 2")
				return
			}
		}
		ch0 <- AssertionFailed.New("test failed")
	}()

	handleMischief(ch1)
}

func handleMischief(ch chan error) {
	err := <-ch
	Panic(Decorate(err, "handle"))
}

func doMischief(ch chan error) {
	defer func() {
		if e := recover(); e != nil {
			err, ok := ErrorFromPanic(e)
			if ok {
				ch <- Decorate(err, "hop 1")
				return
			}
		}
		ch <- AssertionFailed.New("test failed")
	}()

	mischiefAsPanic()
}

func mischiefAsPanic() {
	Panic(mischiefProper())
}

func mischiefProper() error {
	return ExternalError.New("mischief")
}
