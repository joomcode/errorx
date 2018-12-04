package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testNamespace       = NewNamespace("foo")
	testType            = testNamespace.NewType("bar")
	testTypeSilent      = testType.NewSubtype("silent").ApplyModifiers(TypeModifierOmitStackTrace)
	testTypeTransparent = testType.NewSubtype("transparent").ApplyModifiers(TypeModifierTransparent)
	testSubtype0        = testType.NewSubtype("internal")
	testSubtype1        = testSubtype0.NewSubtype("wat")
	testTypeBar1        = testNamespace.NewType("bar1")
	testTypeBar2        = testNamespace.NewType("bar2")
)

func TestError(t *testing.T) {
	err := testType.NewWithNoMessage()
	require.Equal(t, "foo.bar", err.Error())
}

func TestErrorWithMessage(t *testing.T) {
	err := testType.New("oops")
	require.Equal(t, "foo.bar: oops", err.Error())
}

func TestErrorMessageWithCause(t *testing.T) {
	err := testSubtype1.WrapWithNoMessage(testType.New("fatal"))
	require.Equal(t, "foo.bar.internal.wat: foo.bar: fatal", err.Error())
}

func TestErrorWrap(t *testing.T) {
	err0 := testType.NewWithNoMessage()
	err1 := testTypeBar1.Wrap(err0, "a")

	require.Nil(t, Ignore(err1, testTypeBar1))
	require.NotNil(t, Ignore(err1, testType))
}

func TestErrorDecorate(t *testing.T) {
	err0 := testType.NewWithNoMessage()
	err1 := testTypeBar1.Wrap(err0, "a")
	err2 := Decorate(err1, "b")

	require.NotNil(t, Ignore(err2, testTypeBar2))
	require.Nil(t, Ignore(err2, testTypeBar1))
	require.NotNil(t, Ignore(err2, testType))
}

func TestErrorMessages(t *testing.T) {
	t.Run("Subtypes", func(t *testing.T) {
		require.Equal(t, "foo.bar.internal.wat", testSubtype1.NewWithNoMessage().Error())
		require.Equal(t, "foo.bar.internal.wat: oops", testSubtype1.New("oops").Error())
	})

	t.Run("Wrapped", func(t *testing.T) {
		cause := testType.New("poof!")
		require.Equal(t, "foo.bar.internal.wat: foo.bar: poof!", testSubtype1.Wrap(cause, "").Error())
		require.Equal(t, "foo.bar.internal.wat: foo.bar: poof!", testSubtype1.WrapWithNoMessage(cause).Error())
		require.Equal(t, "foo.bar.internal.wat: oops, cause: foo.bar: poof!", testSubtype1.Wrap(cause, "oops").Error())
	})

	t.Run("Complex", func(t *testing.T) {
		innerCause := NewNamespace("c").NewType("d").Wrap(errors.New("Achtung!"), "panic")
		stackedError := testSubtype1.Wrap(testType.Wrap(innerCause, "poof!"), "")
		require.Equal(t, "foo.bar.internal.wat: foo.bar: poof!, cause: c.d: panic, cause: Achtung!", stackedError.Error())
	})
}

func TestImmutableError(t *testing.T) {
	t.Run("Property", func(t *testing.T) {
		err := testType.NewWithNoMessage()
		err1 := err.WithProperty(PropertyPayload(), 1)
		err2 := err1.WithProperty(PropertyPayload(), 2)

		require.True(t, err.errorType.IsOfType(err2.errorType))
		require.Equal(t, err.message, err2.message)

		payload, ok := ExtractPayload(err)
		require.False(t, ok)

		payload, ok = ExtractPayload(err1)
		require.True(t, ok)
		require.EqualValues(t, 1, payload)

		payload, ok = ExtractPayload(err2)
		require.True(t, ok)
		require.EqualValues(t, 2, payload)
	})

	t.Run("Underlying", func(t *testing.T) {
		err := testType.NewWithNoMessage()
		err1 := err.WithUnderlyingErrors(testSubtype0.NewWithNoMessage())
		err2 := err1.WithUnderlyingErrors(testSubtype1.NewWithNoMessage())

		require.True(t, err.errorType.IsOfType(err2.errorType))
		require.Equal(t, err.message, err2.message)

		require.Len(t, err.underlying(), 0)
		require.Len(t, err1.underlying(), 1)
		require.Len(t, err2.underlying(), 2)
	})
}

func TestErrorStackTrace(t *testing.T) {
	err := createErrorFuncInStackTrace(testType)
	output := fmt.Sprintf("%+v", err)
	require.Contains(t, output, "createErrorFuncInStackTrace", output)
	require.Contains(t, output, "TestErrorStackTrace", output)
}

func TestEnhancedStackTrace(t *testing.T) {
	err := createWrappedErrorFuncOuterInStackTrace(testType)
	output := fmt.Sprintf("%+v", err)
	require.Contains(t, output, "createWrappedErrorFuncOuterInStackTrace", output)
	require.Contains(t, output, "createErrorInAnotherGoroutine", output)
}

func TestDecorate(t *testing.T) {
	err := Decorate(testType.NewWithNoMessage(), "ouch!")
	require.Equal(t, "ouch!, cause: foo.bar", err.Error())
	require.True(t, IsOfType(err, testType))
	require.Equal(t, testType, err.Type())
}

func TestUnderlyingInFormat(t *testing.T) {
	err := DecorateMany("this is terribly bad", testTypeBar1.Wrap(testSubtype1.NewWithNoMessage(), "real bad"), testTypeBar2.New("bad"))
	require.Equal(t, "synthetic.wrap: this is terribly bad, cause: foo.bar1: real bad, cause: foo.bar.internal.wat (hidden: foo.bar2: bad)", err.Error())

	err = DecorateMany("this is terribly bad", testTypeBar1.New("real bad"), testTypeBar2.Wrap(testSubtype1.NewWithNoMessage(), "bad"))
	require.Equal(t, "synthetic.wrap: this is terribly bad, cause: foo.bar1: real bad (hidden: foo.bar2: bad, cause: foo.bar.internal.wat)", err.Error())
}

func createErrorFuncInStackTrace(et *Type) *Error {
	err := et.NewWithNoMessage()
	return err
}

func createWrappedErrorFuncOuterInStackTrace(et *Type) *Error {
	return createWrappedErrorFuncInnerInStackTrace(et)
}

func createWrappedErrorFuncInnerInStackTrace(et *Type) *Error {
	channel := make(chan *Error)
	go func() {
		createErrorInAnotherGoroutine(et, channel)
	}()

	errFromChan := <-channel
	return EnhanceStackTrace(errFromChan, "wrap")
}

func createErrorInAnotherGoroutine(et *Type, channel chan *Error) {
	channel <- et.NewWithNoMessage()
}
