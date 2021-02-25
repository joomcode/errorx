// +build go1.13

package errorx

import (
	"errors"
	"fmt"
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
		require.True(t, Cast(unwrapped).Type() == testType)
	})

	t.Run("DecorateForeign", func(t *testing.T) {
		err := Decorate(io.EOF, "")
		unwrapped := errors.Unwrap(err)
		require.NotNil(t, unwrapped)
		require.True(t, errors.Is(unwrapped, io.EOF))
		require.True(t, unwrapped == io.EOF)
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

func TestErrorIs(t *testing.T) {
	t.Run("Trivial", func(t *testing.T) {
		err := testType.NewWithNoMessage()
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
		require.False(t, errors.Is(err, testTypeBar1.NewWithNoMessage()))
	})

	t.Run("Wrap", func(t *testing.T) {
		err := testTypeBar1.Wrap(testType.NewWithNoMessage(), "")
		require.False(t, errors.Is(err, testType.NewWithNoMessage()))
		require.True(t, errors.Is(err, testTypeBar1.NewWithNoMessage()))
	})

	t.Run("Supertype", func(t *testing.T) {
		err := testSubtype0.Wrap(testTypeBar1.NewWithNoMessage(), "")
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
		require.True(t, errors.Is(err, testSubtype0.NewWithNoMessage()))
		require.False(t, errors.Is(err, testTypeBar1.NewWithNoMessage()))
	})

	t.Run("Decorate", func(t *testing.T) {
		err := Decorate(testType.NewWithNoMessage(), "")
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
	})

	t.Run("DecorateForeign", func(t *testing.T) {
		err := Decorate(io.EOF, "")
		require.True(t, errors.Is(err, io.EOF))
	})
}

func TestErrorsAndErrorx(t *testing.T) {
	t.Run("DecoratedForeign", func(t *testing.T) {
		err := fmt.Errorf("error test: %w", testType.NewWithNoMessage())
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
		require.True(t, IsOfType(err, testType))
	})

	t.Run("LayeredDecorate", func(t *testing.T) {
		err := Decorate(fmt.Errorf("error test: %w", testType.NewWithNoMessage()), "test")
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
		require.True(t, IsOfType(err, testType))
	})

	t.Run("LayeredDecorateAgain", func(t *testing.T) {
		err := fmt.Errorf("error test: %w", Decorate(io.EOF, "test"))
		require.True(t, errors.Is(err, io.EOF))
	})

	t.Run("Wrap", func(t *testing.T) {
		err := fmt.Errorf("error test: %w", testType.Wrap(io.EOF, "test"))
		require.False(t, errors.Is(err, io.EOF))
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
	})
}