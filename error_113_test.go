// +build go1.13

package errorx

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorUnwrap(t *testing.T) {
	t.Run("Trivial", func(t *testing.T) {
		err := testType.NewWithNoMessage()
		unwrapped := errors.Unwrap(err)
		require.Nil(t, unwrapped)
	})

	t.Run("Wrap", func(t *testing.T) {
		err := testTypeBar1.Wrap(testType.NewWithNoMessage(), "")
		unwrapped := errors.Unwrap(err)
		require.Nil(t, unwrapped)
	})

	t.Run("WrapForeign", func(t *testing.T) {
		err := testTypeBar1.Wrap(io.EOF, "")
		unwrapped := errors.Unwrap(err)
		require.Nil(t, unwrapped)
	})

	t.Run("Decorate", func(t *testing.T) {
		err := Decorate(testType.NewWithNoMessage(), "")
		unwrapped := errors.Unwrap(err)
		require.NotNil(t, unwrapped)
		require.True(t, IsOfType(unwrapped, testType))
	})

	t.Run("DecorateForeign", func(t *testing.T) {
		err := Decorate(io.EOF, "")
		unwrapped := errors.Unwrap(err)
		require.NotNil(t, unwrapped)
		require.True(t, errors.Is(unwrapped, io.EOF))
	})

	t.Run("Nested", func(t *testing.T) {
		err := Decorate(Decorate(testType.NewWithNoMessage(), ""), "")
		unwrapped := errors.Unwrap(err)
		require.NotNil(t, unwrapped)
		unwrapped = errors.Unwrap(unwrapped)
		require.NotNil(t, unwrapped)
		require.True(t, IsOfType(unwrapped, testType))
	})

	t.Run("NestedWrapped", func(t *testing.T) {
		err := Decorate(testTypeBar1.Wrap(testType.NewWithNoMessage(), ""), "")
		unwrapped := errors.Unwrap(err)
		require.NotNil(t, unwrapped)
		require.True(t, IsOfType(unwrapped, testTypeBar1))
		unwrapped = errors.Unwrap(unwrapped)
		require.Nil(t, unwrapped)
	})

	t.Run("NestedForeign", func(t *testing.T) {
		err := Decorate(Decorate(io.EOF, ""), "")
		unwrapped := errors.Unwrap(err)
		require.NotNil(t, unwrapped)
		unwrapped = errors.Unwrap(unwrapped)
		require.NotNil(t, unwrapped)
		require.True(t, errors.Is(unwrapped, io.EOF))
	})
}
