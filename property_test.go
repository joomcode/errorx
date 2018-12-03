package errorx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoProperty(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		err := testType.New("test")
		property, ok := err.Property(PropertyPayload())
		require.False(t, ok)
		require.Nil(t, property)
	})

	t.Run("Decorated", func(t *testing.T) {
		err := testType.New("test")
		err = Decorate(err, "oops")
		property, ok := err.Property(PropertyPayload())
		require.False(t, ok)
		require.Nil(t, property)
	})

	t.Run("Helper", func(t *testing.T) {
		err := testType.New("test")
		property, ok := ExtractPayload(err)
		require.False(t, ok)
		require.Nil(t, property)
	})
}

var testProperty0 = RegisterProperty("test0")
var testProperty1 = RegisterProperty("test1")

func TestProperty(t *testing.T) {
	t.Run("Different", func(t *testing.T) {
		err := testType.New("test").WithProperty(testProperty0, 42)

		property0, ok := err.Property(testProperty0)
		require.True(t, ok)
		require.EqualValues(t, 42, property0)

		property1, ok := err.Property(testProperty1)
		require.False(t, ok)
		require.Nil(t, property1)
	})

	t.Run("Wrapped", func(t *testing.T) {
		err := testType.New("test").WithProperty(testProperty0, 42)
		err = Decorate(err, "oops")
		err = testTypeBar1.Wrap(err, "wrapped")

		property0, ok := err.Property(testProperty0)
		require.False(t, ok)
		require.Nil(t, property0)

		property1, ok := err.Property(testProperty1)
		require.False(t, ok)
		require.Nil(t, property1)
	})

	t.Run("Decorated", func(t *testing.T) {
		err := testType.New("test").WithProperty(testProperty0, 42)
		err = Decorate(err, "oops")
		err = Decorate(err, "bad")

		property0, ok := err.Property(testProperty0)
		require.True(t, ok)
		require.EqualValues(t, 42, property0)

		property1, ok := err.Property(testProperty1)
		require.False(t, ok)
		require.Nil(t, property1)
	})

	t.Run("FromCause", func(t *testing.T) {
		err := testType.New("test").WithProperty(testProperty0, 42)
		err = Decorate(err, "oops")
		err = Decorate(err, "bad").WithProperty(testProperty1, "-1")

		property0, ok := err.Property(testProperty0)
		require.True(t, ok)
		require.EqualValues(t, 42, property0)

		property1, ok := err.Property(testProperty1)
		require.True(t, ok)
		require.EqualValues(t, "-1", property1)
	})

	t.Run("OverrideCause", func(t *testing.T) {
		err := testType.New("test").WithProperty(testProperty0, 42)
		err = Decorate(err, "oops")

		property0, ok := err.Property(testProperty0)
		require.True(t, ok)
		require.EqualValues(t, 42, property0)

		err = Decorate(err, "bad").WithProperty(testProperty0, "-1")

		property0, ok = err.Property(testProperty0)
		require.True(t, ok)
		require.EqualValues(t, "-1", property0)

		property1, ok := err.Property(testProperty1)
		require.False(t, ok)
		require.Nil(t, property1)
	})
}

func BenchmarkAllocProperty(b *testing.B) {
	const N = 9
	var properties = []Property{}
	for j := 0; j < N; j++ {
		n := fmt.Sprintf("props%d", j)
		properties = append(properties, RegisterProperty(n))
		b.Run(n, func(b *testing.B) {
			for k := 0; k < b.N; k++ {
				err := testTypeSilent.New("test")
				for i := 0; i < j; i++ {
					err = err.WithProperty(properties[i], 42)
				}
			}
		})
	}
}

var sum int

func BenchmarkGetProperty(b *testing.B) {
	const N = 9
	var properties = []Property{}
	for j := 0; j < N; j++ {
		n := fmt.Sprintf("props%d", j)
		properties = append(properties, RegisterProperty(n))
		b.Run(n, func(b *testing.B) {
			err := testTypeSilent.New("test")
			for i := 0; i < j; i++ {
				err = err.WithProperty(properties[i], 42)
			}
			for k := 0; k < b.N; k++ {
				v, ok := err.Property(testProperty0)
				if ok {
					sum += v.(int)
				}
				v, ok = err.Property(properties[j])
				if ok {
					sum += v.(int)
				}
				v, ok = err.Property(properties[0])
				if ok {
					sum += v.(int)
				}
			}
		})
	}
}
