package errorx

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStackTraceStart(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		err := AssertionFailed.New("achtung")
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "New()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("NewWithNoMessage", func(t *testing.T) {
		err := AssertionFailed.NewWithNoMessage()
		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "NewWithNoMessage()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("Wrap", func(t *testing.T) {
		err := AssertionFailed.Wrap(TimeoutElapsed.NewWithNoMessage(), "achtung")
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "Wrap()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("WrapWithNoMessage", func(t *testing.T) {
		err := AssertionFailed.WrapWithNoMessage(TimeoutElapsed.NewWithNoMessage())
		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "WrapWithNoMessage()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("WrapAddStackTrace", func(t *testing.T) {
		err := testTypeSilent.NewWithNoMessage()
		output := fmt.Sprintf("%+v", err)
		require.NotContains(t, output, "TestStackTraceStart", output)

		err = AssertionFailed.Wrap(err, "achtung")
		output = fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "Wrap()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("EnhanceStackTrace", func(t *testing.T) {
		err := EnhanceStackTrace(AssertionFailed.New("achtung"), "enhance")
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "EnhanceStackTrace()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("EnhanceStackTraceWithRaw", func(t *testing.T) {
		err := EnhanceStackTrace(errors.New("achtung"), "enhance")
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "EnhanceStackTrace()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("Decorate", func(t *testing.T) {
		err := Decorate(AssertionFailed.New("achtung"), "enhance")
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "Decorate()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("DecorateWithRaw", func(t *testing.T) {
		err := Decorate(errors.New("achtung"), "enhance")
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "Decorate()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

	t.Run("Raw", func(t *testing.T) {
		err := EnsureStackTrace(errors.New("achtung"))
		output := fmt.Sprintf("%+v", err)
		require.Contains(t, output, "achtung", output)
		require.NotContains(t, output, "EnsureStackTrace()", output)
		require.Contains(t, output, "TestStackTraceStart", output)
	})

}

func TestStackTraceEnhance(t *testing.T) {
	err := stackTestStart()
	output := fmt.Sprintf("%+v", err)

	expected := map[string]int{
		"TestStackTraceEnhance()": 0,
		"stackTestStart()":        0,
		"stackTestWithChan()":     0,
		"stackTest2()":            0,
	}
	checkStackTrace(t, output, expected)
}

func stackTestStart() error {
	ch := make(chan error)
	go stackTestWithChan(ch)
	return EnhanceStackTrace(<-ch, "")
}

func stackTestWithChan(ch chan error) {
	err := stackTest0()
	ch <- err
}

func stackTest0() error {
	return stackTest1()
}

func stackTest1() error {
	return stackTest2()
}

func stackTest2() error {
	return AssertionFailed.New("here be dragons")
}

func TestStackTraceDuplicate(t *testing.T) {
	err := stackTestDuplicate1()
	output := fmt.Sprintf("%+v", err)

	expected := map[string]int{
		"TestStackTraceDuplicate()": 0,
		"stackTestDuplicate1()":     0,
		"stackTestStart1()":         0,
		"stackTestWithChan1()":      0,
		"stackTest21()":             0,
	}
	checkStackTrace(t, output, expected)
}

func stackTestDuplicate1() error {
	return EnhanceStackTrace(stackTestStart1(), "")
}

func stackTestStart1() error {
	ch := make(chan error)
	go stackTestWithChan1(ch)
	return EnhanceStackTrace(<-ch, "")
}

func stackTestWithChan1(ch chan error) {
	err := stackTest11()
	ch <- err
}

func stackTest11() error {
	return stackTest21()
}

func stackTest21() error {
	return AssertionFailed.New("here be dragons")
}

func TestStackTraceDuplicateWithIntermittentFrames(t *testing.T) {
	err := stackTestDuplicate2()
	output := fmt.Sprintf("%+v", err)

	expected := map[string]int{
		"TestStackTraceDuplicateWithIntermittentFrames()": 0,
		"stackTestDuplicate2()":                           0,
		"stackTestStart2()":                               0,
		"enhanceFunc2()":                                  0,
		"stackTestWithChan2()":                            0,
		"stackTest22()":                                   0,
	}
	checkStackTrace(t, output, expected)
}

func stackTestDuplicate2() error {
	err := stackTestStart2()
	return enhanceFunc2(err)
}

func enhanceFunc2(err error) error {
	return EnhanceStackTrace(err, "")
}

func stackTestStart2() error {
	ch := make(chan error)
	go stackTestWithChan2(ch)
	err := <-ch
	return EnhanceStackTrace(err, "")
}

func stackTestWithChan2(ch chan error) {
	err := stackTest12()
	ch <- err
}

func stackTest12() error {
	return stackTest22()
}

func stackTest22() error {
	return AssertionFailed.New("here be dragons and dungeons, too")
}

func checkStackTrace(t *testing.T, output string, expected map[string]int) {
	readByLine(t, output, func(line string) {
		for key := range expected {
			if strings.HasSuffix(line, key) {
				expected[key]++
			}
		}
	})

	for key, value := range expected {
		require.EqualValues(t, 1, value, "Wrong count (%d) of '%s' in:\n%s", value, key, output)
	}
}

func readByLine(t *testing.T, output string, f func(string)) {
	reader := bufio.NewReader(bytes.NewReader([]byte(output)))
	for {
		lineBytes, _, readErr := reader.ReadLine()
		if readErr == io.EOF {
			break
		}

		require.NoError(t, readErr)
		f(string(lineBytes))
	}
}
