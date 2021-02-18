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

func TestErrorIs(t *testing.T) {
	t.Run("Trivial", func(t *testing.T) {
		err := testType.NewWithNoMessage()
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
		require.False(t, errors.Is(err, testTypeBar1.NewWithNoMessage()))
	})

	t.Run("Wrap", func(t *testing.T) {
		err := testTypeBar1.Wrap(testType.NewWithNoMessage(),"")
		require.False(t, errors.Is(err, testType.NewWithNoMessage()))
		require.True(t, errors.Is(err, testTypeBar1.NewWithNoMessage()))
	})

	t.Run("Supertype", func(t *testing.T) {
		err := testSubtype0.Wrap(testTypeBar1.NewWithNoMessage(),"")
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
		require.True(t, errors.Is(err, testSubtype0.NewWithNoMessage()))
		require.False(t, errors.Is(err, testTypeBar1.NewWithNoMessage()))
	})

	t.Run("Decorate", func(t *testing.T) {
		err := Decorate(testType.NewWithNoMessage(),"")
		require.True(t, errors.Is(err, testType.NewWithNoMessage()))
	})

	t.Run("DecorateForeign", func(t *testing.T) {
		err := Decorate(io.EOF,"")
		require.True(t, errors.Is(err, io.EOF))
	})
}

func TestErrorAs(t *testing.T) {
	t.Run("Trivial", func(t *testing.T) {
		err := fooReturnsError()
		target := testType.NewWithNoMessage()
		require.True(t, errors.As(err, &target))
		require.EqualValues(t, "whoops", target.Message())
		output := fmt.Sprintf("%+v", target)
		require.Contains(t, output, "fooReturnsError", output)
	})

	t.Run("TypeChecker", func(t *testing.T) {
		err := fooReturnsError()
		target := testType.ForTypeCheck()
		require.True(t, errors.As(err, &target))
		require.EqualValues(t, "whoops", target.AsError().Message())
		output := fmt.Sprintf("%+v", target)
		require.Contains(t, output, "fooReturnsError", output)
	})

	// Current errors.As allows no customization in this behaviour; if go types are assignable, here we go
	t.Run("NegativeBroken", func(t *testing.T) {
		err := fooReturnsError()
		target := testTypeBar1.NewWithNoMessage()
		require.True(t, errors.As(err, &target))
		require.EqualValues(t, "whoops", target.Message())
		require.True(t, IsOfType(target, testType))
		require.False(t, IsOfType(target, testTypeBar1))
	})

	t.Run("NegativeFixed", func(t *testing.T) {
		err := fooReturnsError()
		target := testTypeBar1.ForTypeCheck()
		require.False(t, errors.As(err, &target))
	})

	t.Run("Negative", func(t *testing.T) {
		err := io.EOF
		target := testTypeBar1.NewWithNoMessage()
		require.False(t, errors.As(err, &target))
	})

	t.Run("Wrap", func(t *testing.T) {
		err := testTypeBar2.Wrap(testTypeBar1.NewWithNoMessage(), "")
		target := testTypeBar2.ForTypeCheck()
		require.True(t, errors.As(err, &target))
		targetBad := testTypeBar1.ForTypeCheck()
		require.False(t, errors.As(err, &targetBad))
	})

	t.Run("Decorate", func(t *testing.T) {
		err := Decorate(testTypeBar1.NewWithNoMessage(), "")
		target := testTypeBar1.ForTypeCheck()
		require.True(t, errors.As(err, &target))
	})

	t.Run("DecorateForeign", func(t *testing.T) {
		err := Decorate(myErr("test"),"")
		var target myErr
		require.True(t, errors.As(err, &target))
		require.EqualValues(t, "test", target.Error())
	})

	t.Run("Parent", func(t *testing.T) {
		err := Decorate(fooReturnsSubtypeError(), "")
		target := testSubtype0.ForTypeCheck()
		require.True(t, errors.As(err, &target))
		targetParent := testType.ForTypeCheck()
		require.True(t, errors.As(err, &targetParent))
		require.EqualValues(t, "whoops", targetParent.AsError().Message())
		output := fmt.Sprintf("%+v", targetParent)
		require.Contains(t, output, "fooReturnsSubtypeError", output)
	})
}

func fooReturnsError() error {
	return testType.New("whoops")
}

func fooReturnsSubtypeError() error {
	return testSubtype0.New("whoops")
}

type myErr string

func (e myErr) Error() string {
	return string(e)
}
