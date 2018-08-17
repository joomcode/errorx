package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeName(t *testing.T) {
	require.Equal(t, "foo.bar", testType.FullName())
}

func TestSubTypeName(t *testing.T) {
	require.Equal(t, "foo.bar.internal.wat", testSubtype1.FullName())
}

func TestRootNamespace(t *testing.T) {
	require.Equal(t, testNamespace, testType.NewWithNoMessage().Type().RootNamespace())
}

func TestSubTypeNamespace(t *testing.T) {
	require.Equal(t, "foo", testSubtype1.RootNamespace().FullName())
}

func TestErrorTypeCheck(t *testing.T) {
	require.True(t, testSubtype1.IsOfType(testSubtype1))
	require.False(t, testSubtype1.IsOfType(NewNamespace("a").NewType("b")))
}

func TestErrorTypeCheckNonErrorx(t *testing.T) {
	require.False(t, IsOfType(errors.New("test"), testSubtype1))
}

func TestErrorTypeUpCast(t *testing.T) {
	require.True(t, testSubtype1.IsOfType(testSubtype0))
	require.True(t, testSubtype1.IsOfType(testType))
}

func TestErrorTypeDownCast(t *testing.T) {
	require.False(t, testSubtype0.IsOfType(testSubtype1))
	require.False(t, testType.IsOfType(testSubtype1))
}

func TestErrorTypeSiblingsCast(t *testing.T) {
	subtype10 := testSubtype0.NewSubtype("wat!")
	subtype11 := testSubtype0.NewSubtype("oops")

	require.False(t, subtype10.IsOfType(subtype11))
	require.False(t, subtype11.IsOfType(subtype10))
}
