package errorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureStackTrace(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		err := EnsureStackTrace(testType.New("good"))
		require.True(t, IsOfType(err, testType))
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "good", output)
		require.Contains(t, output, "TestEnsureStackTrace", output)
	})

	t.Run("NoTrace", func(t *testing.T) {
		err := EnsureStackTrace(testTypeSilent.New("average"))
		require.True(t, IsOfType(err, testType))
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "average", output)
		require.Contains(t, output, "TestEnsureStackTrace", output)
	})

	t.Run("Raw", func(t *testing.T) {
		err := EnsureStackTrace(errors.New("bad"))
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "bad", output)
		require.Contains(t, output, "TestEnsureStackTrace", output)
	})
}

func TestDecorateMany(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		err := DecorateMany("ouch!", testType.NewWithNoMessage())
		require.Equal(t, "ouch!, cause: foo.bar", err.Error())
		require.True(t, IsOfType(err, testType))
		require.Equal(t, testType, err.(*Error).Type())
	})

	t.Run("SingleEmpty", func(t *testing.T) {
		require.Nil(t, DecorateMany("ouch!", nil))
	})

	t.Run("ManyEmpty", func(t *testing.T) {
		require.Nil(t, DecorateMany("ouch!", nil, nil))
		require.Nil(t, DecorateMany("ouch!", nil, nil, nil))
	})

	t.Run("ManySame", func(t *testing.T) {
		err := DecorateMany("ouch!", testType.NewWithNoMessage(), nil, testType.New("bad"))
		require.Equal(t, "ouch!, cause: foo.bar (hidden: foo.bar: bad)", err.Error())
		require.True(t, IsOfType(err, testType))
		require.Equal(t, testType, err.(*Error).Type())
	})

	t.Run("ManyDifferent", func(t *testing.T) {
		err := DecorateMany("ouch!", testTypeBar1.NewWithNoMessage(), testTypeBar2.New("bad"), nil)
		require.Equal(t, "synthetic.wrap: ouch!, cause: foo.bar1 (hidden: foo.bar2: bad)", err.Error())
		require.False(t, IsOfType(err, testTypeBar1))
		require.False(t, IsOfType(err, testTypeBar2))
		require.NotEqual(t, testTypeBar1, err.(*Error).Type())
		require.NotEqual(t, testTypeBar2, err.(*Error).Type())
	})
}
