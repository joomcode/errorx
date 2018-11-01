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

func ExampleError_Format() {
	err := nestedCall()

	simpleOutput := fmt.Sprintf("Error short: %v\n", err)
	verboseOutput := fmt.Sprintf("Error full: %+v", err)

	fmt.Println(simpleOutput)
	fmt.Println(verboseOutput)

	// Example output:
	//Error short: common.assertion_failed: example
	//
	//Error full: common.assertion_failed: example
	// at github.com/joomcode/errorx_test.someFunc()
	//	/Users/username/go/src/github.com/joomcode/errorx/example_test.go:102
	// at github.com/joomcode/errorx_test.nestedCall()
	//	/Users/username/go/src/github.com/joomcode/errorx/example_test.go:98
	// at github.com/joomcode/errorx_test.ExampleError_Format()
	//	/Users/username/go/src/github.com/joomcode/errorx/example_test.go:66
	// <...> more
}

func ExampleEnhanceStackTrace() {
	errCh := make(chan error)
	go func() {
		errCh <- nestedCall()
	}()

	err := <-errCh
	verboseOutput := fmt.Sprintf("Error full: %+v", errorx.EnhanceStackTrace(err, "another goroutine"))
	fmt.Println(verboseOutput)

	// Example output:
	//Error full: another goroutine, cause: common.assertion_failed: example
	// at github.com/joomcode/errorx_test.ExampleEnhanceStackTrace()
	//	/Users/username/go/src/github.com/joomcode/errorx/example_test.go:94
	// at testing.runExample()
	//	/usr/local/Cellar/go/1.10.3/libexec/src/testing/example.go:122
	// at testing.runExamples()
	//	/usr/local/Cellar/go/1.10.3/libexec/src/testing/example.go:46
	// at testing.(*M).Run()
	//	/usr/local/Cellar/go/1.10.3/libexec/src/testing/testing.go:979
	// at main.main()
	//	_testmain.go:146
	// ...
	// (1 duplicated frames)
	// ----------------------------------
	// at github.com/joomcode/errorx_test.someFunc()
	//	/Users/username/go/src/github.com/joomcode/errorx/example_test.go:106
	// at github.com/joomcode/errorx_test.nestedCall()
	//	/Users/username/go/src/github.com/joomcode/errorx/example_test.go:102
	// at github.com/joomcode/errorx_test.ExampleEnhanceStackTrace.func1()
	//	/Users/username/go/src/github.com/joomcode/errorx/example_test.go:90
	// at runtime.goexit()
	//	/usr/local/Cellar/go/1.10.3/libexec/src/runtime/asm_amd64.s:2361
}

func ExampleIgnore() {
	err := errorx.IllegalArgument.NewWithNoMessage()
	err = errorx.Decorate(err, "more info")

	fmt.Println(err)
	fmt.Println(errorx.Ignore(err, errorx.IllegalArgument))
	fmt.Println(errorx.Ignore(err, errorx.AssertionFailed))

	// Output:
	// more info, cause: common.illegal_argument
	// <nil>
	// more info, cause: common.illegal_argument
}

func ExampleIgnoreWithTrait() {
	err := errorx.TimeoutElapsed.NewWithNoMessage()
	err = errorx.Decorate(err, "more info")

	fmt.Println(err)
	fmt.Println(errorx.IgnoreWithTrait(err, errorx.Timeout()))
	fmt.Println(errorx.IgnoreWithTrait(err, errorx.NotFound()))

	// Output:
	// more info, cause: common.timeout
	// <nil>
	// more info, cause: common.timeout
}

func ExampleIsOfType() {
	err0 := errorx.DataUnavailable.NewWithNoMessage()
	err1 := errorx.Decorate(err0, "decorated")
	err2 := errorx.RejectedOperation.Wrap(err0, "wrapped")

	fmt.Println(errorx.IsOfType(err0, errorx.DataUnavailable))
	fmt.Println(errorx.IsOfType(err1, errorx.DataUnavailable))
	fmt.Println(errorx.IsOfType(err2, errorx.DataUnavailable))

	// Output:
	// true
	// true
	// false
}

func ExampleTypeSwitch() {
	err := errorx.DataUnavailable.NewWithNoMessage()

	switch errorx.TypeSwitch(err, errorx.DataUnavailable) {
	case errorx.DataUnavailable:
		fmt.Println("good")
	case nil:
		fmt.Println("bad")
	default:
		fmt.Println("bad")
	}

	switch errorx.TypeSwitch(nil, errorx.DataUnavailable) {
	case errorx.DataUnavailable:
		fmt.Println("bad")
	case nil:
		fmt.Println("good")
	default:
		fmt.Println("bad")
	}

	switch errorx.TypeSwitch(err, errorx.TimeoutElapsed) {
	case errorx.TimeoutElapsed:
		fmt.Println("bad")
	case nil:
		fmt.Println("bad")
	default:
		fmt.Println("good")
	}

	// Output:
	// good
	// good
	// good
}

func ExampleTraitSwitch() {
	err := errorx.TimeoutElapsed.NewWithNoMessage()

	switch errorx.TraitSwitch(err, errorx.Timeout()) {
	case errorx.Timeout():
		fmt.Println("good")
	case errorx.CaseNoError():
		fmt.Println("bad")
	default:
		fmt.Println("bad")
	}

	switch errorx.TraitSwitch(nil, errorx.Timeout()) {
	case errorx.Timeout():
		fmt.Println("bad")
	case errorx.CaseNoError():
		fmt.Println("good")
	default:
		fmt.Println("bad")
	}

	switch errorx.TraitSwitch(err, errorx.NotFound()) {
	case errorx.NotFound():
		fmt.Println("bad")
	case errorx.CaseNoError():
		fmt.Println("bad")
	default:
		fmt.Println("good")
	}

	// Output:
	// good
	// good
	// good
}

func nestedCall() error {
	return someFunc()
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
