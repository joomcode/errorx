package errorx

import (
	"errors"
	"fmt"
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

func testBuilderRespectsNoStackTraceMarkerFrame() error {
	return testType.NewWithNoMessage()
}

func TestBuilderRespectsNoStackTrace(t *testing.T) {
	wrapperErrorTypes := []*Type{testTypeSilent, testTypeSilentTransparent}

	for _, et := range wrapperErrorTypes {
		t.Run(et.String(), func(t *testing.T) {
			t.Run("Naked", func(t *testing.T) {
				err := NewErrorBuilder(et).
					WithCause(errors.New("naked error")).
					Create()
				require.Nil(t, err.stackTrace)
			})

			t.Run("WithoutStacktrace", func(t *testing.T) {
				err := NewErrorBuilder(et).
					WithCause(testTypeSilent.NewWithNoMessage()).
					Create()
				require.Nil(t, err.stackTrace)
			})

			t.Run("WithStacktrace", func(t *testing.T) {
				cause := testBuilderRespectsNoStackTraceMarkerFrame()
				err := NewErrorBuilder(et).
					WithCause(cause).
					Create()
				require.Same(t, err.stackTrace, Cast(cause).stackTrace)
				require.Contains(t, fmt.Sprintf("%+v", err), "testBuilderRespectsNoStackTraceMarkerFrame")
			})
		})
	}
}
