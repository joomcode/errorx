package errorx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testTrait0 = RegisterTrait("test0")
	testTrait1 = RegisterTrait("test1")
	testTrait2 = RegisterTrait("test2")

	traitTestNamespace             = NewNamespace("traits")
	traitTestNamespace2            = NewNamespace("traits2", testTrait0)
	traitTestNamespace2Child       = traitTestNamespace2.NewSubNamespace("child", testTrait1)
	traitTestError                 = traitTestNamespace.NewType("simple", testTrait1)
	traitTestError2                = traitTestNamespace2.NewType("simple", testTrait2)
	traitTestError3                = traitTestNamespace2Child.NewType("simple", testTrait2)
	traitTestTimeoutError          = traitTestNamespace.NewType("timeout", Timeout())
	traitTestTemporaryTimeoutError = traitTestTimeoutError.NewSubtype("temporary", Temporary())
)

func TestTrait(t *testing.T) {
	t.Run("Negative", func(t *testing.T) {
		err := traitTestError.New("test")
		require.False(t, IsTemporary(err))
	})

	t.Run("Positive", func(t *testing.T) {
		err := traitTestError.New("test")
		require.True(t, HasTrait(err, testTrait1))
	})

	t.Run("SubType", func(t *testing.T) {
		err := traitTestTimeoutError.New("test")
		require.True(t, IsTimeout(err))
		require.False(t, IsTemporary(err))

		err = traitTestTemporaryTimeoutError.New("test")
		require.True(t, IsTimeout(err))
		require.True(t, IsTemporary(err))
	})

	t.Run("Wrap", func(t *testing.T) {
		err := traitTestTimeoutError.New("test")
		err = traitTestError2.Wrap(err, "")

		require.False(t, IsTimeout(err))
		require.True(t, HasTrait(err, testTrait0))
		require.False(t, HasTrait(err, testTrait1))
		require.True(t, HasTrait(err, testTrait2))
	})

	t.Run("Decorate", func(t *testing.T) {
		err := traitTestTimeoutError.New("test")
		err = Decorate(err, "")

		require.True(t, IsTimeout(err))
		require.False(t, IsTemporary(err))
	})
}

func TestTraitNamespace(t *testing.T) {
	t.Run("Negative", func(t *testing.T) {
		err := traitTestError.New("test")
		require.False(t, HasTrait(err, testTrait0))
		require.True(t, HasTrait(err, testTrait1))
		require.False(t, HasTrait(err, testTrait2))
	})

	t.Run("Inheritance", func(t *testing.T) {
		err := traitTestError2.New("test")
		require.True(t, HasTrait(err, testTrait0))
		require.False(t, HasTrait(err, testTrait1))
		require.True(t, HasTrait(err, testTrait2))
	})

	t.Run("DoubleInheritance", func(t *testing.T) {
		err := traitTestError3.New("test")
		require.True(t, HasTrait(err, testTrait0))
		require.True(t, HasTrait(err, testTrait1))
		require.True(t, HasTrait(err, testTrait2))
	})
}
