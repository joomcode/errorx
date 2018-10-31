package errorx

import (
	"fmt"
	"io"
	"strings"
)

// Error is an instance of error object.
// At the moment of creation, Error collects information based on context, creation modifiers and type it belongs to.
// Error is mostly immutable, and distinct errors composition is achieved through wrap.
type Error struct {
	message     string
	errorType   *Type
	cause       error
	underlying  []error
	stackTrace  *stackTrace
	transparent bool
	properties  map[Property]interface{}
}

var _ fmt.Formatter = (*Error)(nil)

// WithProperty adds a dynamic property to error instance.
// If an error already contained another value for the same property, it is overwritten.
// It is a caller's responsibility to accumulate and update a property, if needed.
// Dynamic properties is a brittle mechanism and should therefore be used with care and in a simple and robust manner.
func (e *Error) WithProperty(key Property, value interface{}) *Error {
	if e.properties == nil {
		e.properties = make(map[Property]interface{}, 1)
	}

	e.properties[key] = value
	return e
}

// WithUnderlyingErrors adds multiple additional related (hidden, suppressed) errors to be used exclusively in error output.
// Note that these errors make no other effect whatsoever: their traits, types, properties etc. are lost on the observer.
// Consider using errorx.DecorateMany instead.
func (e *Error) WithUnderlyingErrors(errs ...error) *Error {
	for _, err := range errs {
		if err == nil {
			continue
		}

		e.underlying = append(e.underlying, err)
	}

	return e
}

// Property extracts a dynamic property value from an error.
// A property may belong to this error or be extracted from the original cause.
// The transparency rules are respected to some extent: both the original cause and the transparent wrapper
// may have accessible properties, but an opaque wrapper hides the original properties.
func (e *Error) Property(key Property) (interface{}, bool) {
	cause := e
	for cause != nil {
		value, ok := cause.properties[key]
		if ok {
			return value, true
		}

		if !cause.transparent {
			break
		}

		cause = Cast(cause.Cause())
	}

	return nil, false
}

// HasTrait checks if an error possesses the expected trait.
// Trait check works just as a type check would: opaque wrap hides the traits of the cause.
// Traits are always properties of a type rather than of an instance, so trait check is an alternative to a type check.
// This alternative is preferable, though, as it is less brittle and generally creates less of a dependency.
func (e *Error) HasTrait(key Trait) bool {
	cause := e
	for cause != nil {
		if !cause.transparent {
			_, ok := cause.errorType.traits[key]
			return ok
		}

		cause = Cast(cause.Cause())
	}

	return false
}

// IsOfType is a proper type check for an error.
// It takes the transparency and error types hierarchy into account,
// so that type check against any supertype of the original cause passes.
func (e *Error) IsOfType(t *Type) bool {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.errorType.IsOfType(t)
		}

		cause = Cast(cause.Cause())
	}

	return false
}

// Type returns the exact type of this error.
// With transparent wrapping, such as in Decorate(), returns the type of the original cause.
// The result is always not nil, even if the resulting type is impossible to successfully type check against.
//
// NB: the exact error type may fail an equality check where a IsOfType() check would succeed.
// This may happen if a type is checked against one of its supertypes, for example.
// Therefore, handle direct type checks with care or avoid it altogether and use TypeSwitch() or IsForType() instead.
func (e *Error) Type() *Type {
	cause := e
	for cause != nil {
		if !cause.transparent {
			return cause.errorType
		}

		cause = Cast(cause.Cause())
	}

	return foreignType
}

// Message returns a message of this particular error, disregarding the cause.
// The result of this method, like a result of an Error() method, should never be used to infer the meaning of an error.
// In most cases, message is only used as a part of formatting to print error contents into a log file.
// Manual extraction may be required, however, to transform an error into another format - say, API response.
func (e *Error) Message() string {
	return e.message
}

// Cause returns the immediate (wrapped) cause of current error.
// This method could be used to dig for root cause of the error, but it is not advised to do so.
// Errors should not require a complex navigation through causes to be properly handled, and the need to do so is a code smell.
// Manually extracting cause defeats features such as opaque wrap, behaviour of properties etc.
// This method is, therefore, reserved for system utilities, not for general use.
func (e *Error) Cause() error {
	return e.cause
}

// Format implements the Formatter interface.
// Supported verbs:
//
// 		%s		simple message output
// 		%v		same as %s
// 		%+v		full output complete with a stack trace
//
// In is nearly always preferable to use %+v format.
// If a stack trace is not required, it should be omitted at the moment of creation rather in formatting.
func (e *Error) Format(s fmt.State, verb rune) {
	message := e.fullMessage()
	switch verb {
	case 'v':
		io.WriteString(s, message)
		if s.Flag('+') {
			e.stackTrace.Format(s, verb)
		}
	case 's':
		io.WriteString(s, message)
	}
}

// Error implements the error interface.
// A result is the same as with %s formatter and does not contain a stack trace.
func (e *Error) Error() string {
	return e.fullMessage()
}

func (e *Error) fullMessage() string {
	if e.transparent {
		return e.messageWithUnderlyingInfo()
	} else {
		return joinStringsIfNonEmpty(": ", e.errorType.FullName(), e.messageWithUnderlyingInfo())
	}
}

func (e *Error) messageWithUnderlyingInfo() string {
	return joinStringsIfNonEmpty(" ", e.messageText(), e.underlyingInfo())
}

func (e *Error) underlyingInfo() string {
	if len(e.underlying) == 0 {
		return ""
	}

	infos := make([]string, 0, len(e.underlying))
	for _, err := range e.underlying {
		infos = append(infos, err.Error())
	}

	return fmt.Sprintf("(hidden: %s)", joinStringsIfNonEmpty(", ", infos...))
}

func (e *Error) messageText() string {
	if e.Cause() == nil {
		return e.message
	}

	underlyingFullMessage := e.Cause().Error()
	return joinStringsIfNonEmpty(", cause: ", e.message, underlyingFullMessage)
}

func joinStringsIfNonEmpty(delimiter string, parts ...string) string {
	filteredParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) > 0 {
			filteredParts = append(filteredParts, part)
		}
	}

	return strings.Join(filteredParts, delimiter)
}
