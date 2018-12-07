package errorx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
var testInfoProperty2 = RegisterPrintableProperty("prop2")
var testInfoProperty3 = RegisterPrintableProperty("prop3")

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

func TestPrintableProperty(t *testing.T) {
	err := testTypeSilent.New("test").WithProperty(testInfoProperty2, "hello world")
	t.Run("Simple", func(t *testing.T) {
		assert.Equal(t, "foo.bar.silent: test {prop2: hello world}", err.Error())
	})

	t.Run("Overwrite", func(t *testing.T) {
		err := err.WithProperty(testInfoProperty2, "cruel world")
		assert.Equal(t, "foo.bar.silent: test {prop2: cruel world}", err.Error())
	})

	t.Run("AddMore", func(t *testing.T) {
		err := err.WithProperty(testInfoProperty3, struct{ a int }{1})
		assert.Equal(t, "foo.bar.silent: test {prop3: {1}, prop2: hello world}", err.Error())
	})

	t.Run("NonPrintableIsInvisible", func(t *testing.T) {
		err := err.WithProperty(testProperty0, "nah")
		assert.Equal(t, "foo.bar.silent: test {prop2: hello world}", err.Error())
	})

	t.Run("WithUnderlying", func(t *testing.T) {
		err := err.WithUnderlyingErrors(testTypeSilent.New("underlying"))
		assert.Equal(t, "foo.bar.silent: test {prop2: hello world} (hidden: foo.bar.silent: underlying)", err.Error())
	})

	err2 := Decorate(err, "oops")
	t.Run("Decorate", func(t *testing.T) {
		assert.Equal(t, "oops, cause: foo.bar.silent: test {prop2: hello world}", err2.Error())
	})

	t.Run("DecorateAndAddMore", func(t *testing.T) {
		err := err2.WithProperty(testInfoProperty3, struct{ a int }{1})
		assert.Equal(t, "oops {prop3: {1}}, cause: foo.bar.silent: test {prop2: hello world}", err.Error())
	})

	t.Run("DecorateAndAddSame", func(t *testing.T) {
		err := err2.WithProperty(testInfoProperty2, "cruel world")
		assert.Equal(t, "oops {prop2: cruel world}, cause: foo.bar.silent: test {prop2: hello world}", err.Error())
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
