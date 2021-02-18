package errorx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsWrapper(t *testing.T) {
	require.False(t, reflect.TypeOf(Error{}).AssignableTo(reflect.TypeOf(typeCheckTarget{})))
}