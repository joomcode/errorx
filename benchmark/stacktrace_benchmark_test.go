package benchmark

import (
	"errors"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/joomcode/errorx"
)

var errorSink error

func BenchmarkSimpleError10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errorSink = function0(10, createSimpleError)
	}
	consumeResult(errorSink)
}

func BenchmarkErrorxError10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errorSink = function0(10, createSimpleErrorxError)
	}
	consumeResult(errorSink)
}

func BenchmarkStackTraceErrorxError10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errorSink = function0(10, createErrorxError)
	}
	consumeResult(errorSink)
}

func BenchmarkSimpleError100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errorSink = function0(100, createSimpleError)
	}
	consumeResult(errorSink)
}

func BenchmarkErrorxError100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errorSink = function0(100, createSimpleErrorxError)
	}
	consumeResult(errorSink)
}

func BenchmarkStackTraceErrorxError100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errorSink = function0(100, createErrorxError)
	}
	consumeResult(errorSink)
}

func BenchmarkStackTraceNaiveError100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errorSink = function0(100, createNaiveError)
	}
	consumeResult(errorSink)
}

func BenchmarkSimpleErrorPrint100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := function0(100, createSimpleError)
		emulateErrorPrint(err)
		errorSink = err
	}
	consumeResult(errorSink)
}

func BenchmarkErrorxErrorPrint100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := function0(100, createSimpleErrorxError)
		emulateErrorPrint(err)
		errorSink = err
	}
	consumeResult(errorSink)
}

func BenchmarkStackTraceErrorxErrorPrint100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := function0(100, createErrorxError)
		emulateErrorPrint(err)
		errorSink = err
	}
	consumeResult(errorSink)
}

func BenchmarkStackTraceNaiveErrorPrint100(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := function0(100, createNaiveError)
		emulateErrorPrint(err)
		errorSink = err
	}
	consumeResult(errorSink)
}

func createSimpleError() error {
	return errors.New("benchmark")
}

var (
	Errors            = errorx.NewNamespace("errorx.benchmark")
	NoStackTraceError = Errors.NewType("no_stack_trace").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	StackTraceError   = Errors.NewType("stack_trace")
)

func createSimpleErrorxError() error {
	return NoStackTraceError.New("benchmark")
}

func createErrorxError() error {
	return StackTraceError.New("benchmark")
}

type naiveError struct {
	stack []byte
}

func (err naiveError) Error() string {
	return fmt.Sprintf("benchmark\n%s", err.stack)
}

func createNaiveError() error {
	return naiveError{stack: debug.Stack()}
}

func function0(depth int, generate func() error) error {
	if depth == 0 {
		return generate()
	}

	switch depth % 3 {
	case 0:
		return function1(depth-1, generate)
	case 1:
		return function2(depth-1, generate)
	default:
		return function3(depth-1, generate)
	}
}

func function1(depth int, generate func() error) error {
	if depth == 0 {
		return generate()
	}

	return function4(depth-1, generate)
}

func function2(depth int, generate func() error) error {
	if depth == 0 {
		return generate()
	}

	return function4(depth-1, generate)
}

func function3(depth int, generate func() error) error {
	if depth == 0 {
		return generate()
	}

	return function4(depth-1, generate)
}

func function4(depth int, generate func() error) error {
	switch depth {
	case 0:
		return generate()
	default:
		return function0(depth-1, generate)
	}
}

type sinkError struct {
	value int
}

func (sinkError) Error() string {
	return ""
}

// Perform error formatting and consume the result to disallow optimizations against output
func emulateErrorPrint(err error) {
	output := fmt.Sprintf("%+v", err)
	if len(output) > 10000 && output[1000:1004] == "DOOM" {
		panic("this was not supposed to happen")
	}
}

// Consume error with a possible side effect to disallow optimizations against err
func consumeResult(err error) {
	if e, ok := err.(sinkError); ok && e.value == 1 {
		panic("this was not supposed to happen")
	}
}

// A public function to discourage optimizations against errorSink variable
func ExportSink() error {
	return errorSink
}
