// +build go1.13

package errorx

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type otherError struct {
	message string
}

func (e *otherError) Error() string {
	return e.message
}

func TestErrorsIs(t *testing.T) {
	errInner := testType.NewWithNoMessage()

	t.Run("Wrap", func(t *testing.T) {
		err0 := testTypeBar1.Wrap(errInner, "a")

		require.True(t, errors.Is(err0, errInner))
		require.True(t, err0.Is(errInner))

		err1 := testTypeBar1.NewWithNoMessage()
		require.False(t, errors.Is(err1, errInner))
		require.False(t, err1.Is(errInner))

		err2 := testTypeBar2.WrapWithNoMessage(err1)
		require.False(t, errors.Is(err2, errInner))
		require.False(t, err2.Is(errInner))
		require.True(t, errors.Is(err2, err1))
		require.True(t, err2.Is(err1))

		err3 := testTypeTransparent.WrapWithNoMessage(testTypeBar1.NewWithNoMessage())
		err4 := testTypeTransparent.WrapWithNoMessage(testTypeBar2.NewWithNoMessage())
		require.False(t, errors.Is(err3, err4))
		require.False(t, err3.Is(err4))
		require.False(t, errors.Is(err4, err3))
		require.False(t, err4.Is(err3))
	})

	t.Run("Decorate", func(t *testing.T) {
		err1 := testTypeBar1.Wrap(errInner, "a")
		err2 := Decorate(err1, "b")

		require.True(t, errors.Is(err2, err1))
		require.True(t, err2.Is(err1))
		require.True(t, errors.Is(err2, errInner))
		require.True(t, err2.Is(errInner))
	})

	t.Run("Builtins", func(t *testing.T) {
		err0 := testType.WrapWithNoMessage(context.Canceled)
		err1 := testTypeBar1.WrapWithNoMessage(err0)

		require.True(t, errors.Is(err0, context.Canceled))
		require.True(t, err0.Is(context.Canceled))
		require.True(t, errors.Is(err1, context.Canceled))
		require.True(t, err1.Is(context.Canceled))
	})
}

func TestErrorsAs(t *testing.T) {
	t.Run("errorx.Error", func(t *testing.T) {
		errInner := testType.NewWithNoMessage()
		err0 := testTypeBar1.Wrap(errInner, "a")
		err1 := testTypeBar2.Wrap(err0, "b")

		var e0 *Error
		require.True(t, errors.As(err0, &e0))
		require.Equal(t, e0.message, err0.message)
		require.True(t, errors.Is(e0, errInner))

		var e2 *Error
		require.True(t, err0.As(&e2))
		require.Equal(t, e2.message, err0.message)
		require.True(t, errors.Is(e2, errInner))

		var e3 *Error
		require.True(t, errors.As(err1, &e3))
		require.Equal(t, e3.message, err1.message)
		require.True(t, errors.Is(e3, errInner))

		var e4 *Error
		require.True(t, err1.As(&e4))
		require.Equal(t, e4.message, err1.message)
		require.True(t, errors.Is(e4, errInner))
	})

	t.Run("OtherError", func(t *testing.T) {
		errOther := &otherError{message: "foo"}
		err := testType.WrapWithNoMessage(errOther)

		var oe0 *otherError
		require.True(t, errors.As(err, &oe0))
		require.Equal(t, oe0.message, errOther.message)

		var oe1 *otherError
		require.True(t, err.As(&oe1))
		require.Equal(t, oe1.message, errOther.message)
	})
}

func TestErrorUnwrap(t *testing.T) {
	errInner := testType.NewWithNoMessage()

	t.Run("Wrap", func(t *testing.T) {
		err := testTypeBar1.Wrap(errInner, "a")

		unwrapped := errors.Unwrap(err)
		require.True(t, errors.Is(unwrapped, errInner))
		require.False(t, errors.Is(unwrapped, err))

		unwrapped = err.Unwrap()
		require.True(t, errors.Is(unwrapped, errInner))
		require.False(t, errors.Is(unwrapped, err))
	})

	t.Run("Decorate", func(t *testing.T) {
		err0 := testTypeBar1.Wrap(errInner, "a")
		err1 := Decorate(err0, "b")

		unwrapped := errors.Unwrap(err1)
		require.True(t, errors.Is(unwrapped, err0))
		require.True(t, errors.Is(unwrapped, errInner))

		unwrapped = err1.Unwrap()
		require.True(t, errors.Is(unwrapped, err0))
		require.True(t, errors.Is(unwrapped, errInner))
	})

	t.Run("Builtins", func(t *testing.T) {
		err0 := testType.WrapWithNoMessage(context.Canceled)
		err1 := testTypeBar1.WrapWithNoMessage(context.Canceled)
		err2 := testTypeTransparent.WrapWithNoMessage(context.Canceled)
		err3 := testTypeTransparent.WrapWithNoMessage(err1)

		require.Equal(t, context.Canceled, errors.Unwrap(err0))
		require.Equal(t, context.Canceled, err0.Unwrap())
		require.Equal(t, context.Canceled, errors.Unwrap(err1))
		require.Equal(t, context.Canceled, err1.Unwrap())
		require.Equal(t, context.Canceled, errors.Unwrap(err2))
		require.Equal(t, context.Canceled, err2.Unwrap())

		err3Unwrapped := Cast(errors.Unwrap(err3))
		require.NotNil(t, err3Unwrapped)
		require.True(t, err3Unwrapped.IsOfType(err1.Type()))
		require.True(t, errors.Is(err3Unwrapped, context.Canceled))

		err3Unwrapped = Cast(err3.Unwrap())
		require.NotNil(t, err3Unwrapped)
		require.True(t, err3Unwrapped.IsOfType(err1.Type()))
		require.True(t, errors.Is(err3Unwrapped, context.Canceled))
	})
}
