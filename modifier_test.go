package errorx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	modifierTestNamespace                 = NewNamespace("modifier")
	modifierTestNamespaceTransparent      = NewNamespace("modifierTransparent").ApplyModifiers(TypeModifierTransparent)
	modifierTestNamespaceTransparentChild = modifierTestNamespaceTransparent.NewSubNamespace("child")
	modifierTestError                     = modifierTestNamespace.NewType("foo")
	modifierTestErrorNoTrace              = modifierTestNamespace.NewType("bar").ApplyModifiers(TypeModifierOmitStackTrace)
	modifierTestErrorNoTraceChild         = modifierTestErrorNoTrace.NewSubtype("child")
	modifierTestErrorTransparent          = modifierTestNamespaceTransparent.NewType("simple")
	modifierTestErrorGrandchild           = modifierTestNamespaceTransparentChild.NewType("all").ApplyModifiers(TypeModifierOmitStackTrace)
)

func TestTypeModifier(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		err := modifierTestError.New("test")
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "errorx/modifier_test.go")
	})

	t.Run("NoTrace", func(t *testing.T) {
		err := modifierTestErrorNoTrace.New("test")
		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "errorx/modifier_test.go")
	})
}

func TestTypeModifierInheritance(t *testing.T) {
	t.Run("Type", func(t *testing.T) {
		err := modifierTestErrorNoTraceChild.New("test")
		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "errorx/modifier_test.go")
	})

	t.Run("Namespace", func(t *testing.T) {
		err := modifierTestErrorTransparent.Wrap(AssertionFailed.New("test"), "boo")
		require.True(t, err.IsOfType(AssertionFailed))
	})

	t.Run("Deep", func(t *testing.T) {
		err := modifierTestErrorGrandchild.Wrap(AssertionFailed.New("test"), "boo")
		require.True(t, err.IsOfType(AssertionFailed))

		err = modifierTestErrorGrandchild.New("test")
		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "errorx/modifier_test.go")
	})
}
