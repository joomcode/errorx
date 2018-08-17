package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilderTransparency(t *testing.T) {
	t.Run("Raw", func(t *testing.T) {
		err := NewErrorBuilder(testType).WithCause(errors.New("bad thing")).Transparent().Create()
		require.False(t, err.IsOfType(testType))
		require.NotEqual(t, testType, err.Type())
	})

	t.Run("RawWithModifier", func(t *testing.T) {
		err := NewErrorBuilder(testTypeTransparent).WithCause(errors.New("bad thing")).Create()
		require.False(t, err.IsOfType(testType))
		require.NotEqual(t, testType, err.Type())
	})
}
