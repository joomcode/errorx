package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIgnoreWithTrait(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		require.Error(t, IgnoreWithTrait(TimeoutElapsed.NewWithNoMessage()))
	})

	t.Run("AnotherTrait", func(t *testing.T) {
		require.Error(t, IgnoreWithTrait(TimeoutElapsed.NewWithNoMessage(), NotFound()))
	})

	t.Run("Positive", func(t *testing.T) {
		require.NoError(t, IgnoreWithTrait(TimeoutElapsed.NewWithNoMessage(), Timeout()))
	})

	t.Run("OneOfMany", func(t *testing.T) {
		require.NoError(t, IgnoreWithTrait(TimeoutElapsed.NewWithNoMessage(), NotFound(), Timeout()))
	})
}

func TestGetTypeName(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		require.EqualValues(t, "common.assertion_failed", GetTypeName(AssertionFailed.NewWithNoMessage()))
	})

	t.Run("Wrap", func(t *testing.T) {
		require.EqualValues(t, "common.illegal_state", GetTypeName(IllegalState.WrapWithNoMessage(AssertionFailed.NewWithNoMessage())))
	})

	t.Run("Decorate", func(t *testing.T) {
		require.EqualValues(t, "common.assertion_failed", GetTypeName(Decorate(AssertionFailed.NewWithNoMessage(), "")))
	})

	t.Run("Nil", func(t *testing.T) {
		require.EqualValues(t, "", GetTypeName(nil))
	})

	t.Run("Raw", func(t *testing.T) {
		require.EqualValues(t, "", GetTypeName(errors.New("test")))
	})

	t.Run("DecoratedRaw", func(t *testing.T) {
		require.EqualValues(t, "", GetTypeName(Decorate(errors.New("test"), "")))
	})
}
