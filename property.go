package errorx

import (
	"context"
)

// Property is a key to a dynamic property of an error.
// Property value belongs to an error instance only, never inherited from a type.
// Property visibility is hindered by Wrap, preserved by Decorate.
type Property struct {
	id    int64
	label string
}

// RegisterProperty registers a new property key.
// It is used both to add a dynamic property to an error instance, and to extract property value back from error.
func RegisterProperty(label string) Property {
	return newProperty(label)
}

// PropertyContext is a context property, value is expected to be of context.Context type.
func PropertyContext() Property {
	return propertyContext
}

// PropertyPayload is a payload property, value may contain user defined structure with arbitrary data passed along with an error.
func PropertyPayload() Property {
	return propertyPayload
}

// WithContext is a statically typed helper to add a context property to an error.
func WithContext(err *Error, ctx context.Context) *Error {
	return err.WithProperty(PropertyContext(), ctx)
}

// ExtractContext is a statically typed helper to extract a context property from an error.
func ExtractContext(err error) (context.Context, bool) {
	rawCtx, ok := ExtractProperty(err, PropertyContext())
	if !ok {
		return nil, false
	}

	return rawCtx.(context.Context), true
}

// WithPayload is a helper to add a payload property to an error.
func WithPayload(err *Error, payload interface{}) *Error {
	return err.WithProperty(PropertyPayload(), payload)
}

// ExtractPayload is a helper to extract a payload property from an error.
func ExtractPayload(err error) (interface{}, bool) {
	return ExtractProperty(err, PropertyPayload())
}

// ExtractProperty attempts to extract a property value by a provided key.
// A property may belong to this error or be extracted from the original cause.
func ExtractProperty(err error, key Property) (interface{}, bool) {
	typedErr := Cast(err)
	if typedErr == nil {
		return nil, false
	}

	return typedErr.Property(key)
}

var (
	propertyContext = RegisterProperty("ctx")
	propertyPayload = RegisterProperty("payload")
)

func newProperty(label string) Property {
	return Property{
		id:    nextInternalID(),
		label: label,
	}
}
