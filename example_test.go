package errorx_test

import (
	"fmt"
	"github.com/joomcode/errorx"
)

func ExampleDecorate() {
	err := someFunc()
	fmt.Println(err.Error())

	err = errorx.Decorate(err, "decorate")
	fmt.Println(err.Error())

	err = errorx.Decorate(err, "outer decorate")
	fmt.Println(err.Error())

	// Output: common.assertion_failed: example
	// decorate, cause: common.assertion_failed: example
	// outer decorate, cause: decorate, cause: common.assertion_failed: example
}

func ExampleDecorateMany() {
	err0 := someFunc()
	err1 := someFunc()
	err := errorx.DecorateMany("both calls failed", err0, err1)
	fmt.Println(err.Error())

	// Output: both calls failed, cause: common.assertion_failed: example (hidden: common.assertion_failed: example)
}

func ExampleError_WithUnderlyingErrors() {
	fn := func() error {
		bytes, err := getBodyAndError()
		if err != nil {
			_, unmarshalErr := getDetailsFromBody(bytes)
			if unmarshalErr != nil {
				return errorx.AssertionFailed.Wrap(err, "failed to read details").WithUnderlyingErrors(unmarshalErr)
			}
		}

		return nil
	}

	fmt.Println(fn().Error())
	// Output: common.assertion_failed: failed to read details, cause: common.assertion_failed: example (hidden: common.illegal_format)
}

func ExampleType_Wrap() {
	originalErr := errorx.IllegalArgument.NewWithNoMessage()
	err := errorx.AssertionFailed.Wrap(originalErr, "wrapped")

	fmt.Println(errorx.IsOfType(originalErr, errorx.IllegalArgument))
	fmt.Println(errorx.IsOfType(err, errorx.IllegalArgument))
	fmt.Println(errorx.IsOfType(err, errorx.AssertionFailed))
	fmt.Println(err.Error())

	// Output:
	// true
	// false
	// true
	// common.assertion_failed: wrapped, cause: common.illegal_argument
}

func someFunc() error {
	return errorx.AssertionFailed.New("example")
}

func getBodyAndError() ([]byte, error) {
	return nil, errorx.AssertionFailed.New("example")
}

func getDetailsFromBody(s []byte) (string, error) {
	return "", errorx.IllegalFormat.New(string(s))
}
