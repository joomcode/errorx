package errorx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	s := &testSubscriber{}
	RegisterTypeSubscriber(s)

	require.Contains(t, s.namespaces, CommonErrors.Key())
	require.Contains(t, s.types, AssertionFailed)

	ns := NewNamespace("TestRegistry")
	require.Contains(t, s.namespaces, ns.Key())

	errorType := ns.NewType("Test")
	require.Contains(t, s.types, errorType)
}

type testSubscriber struct {
	types      []*Type
	namespaces []NamespaceKey
}

func (s *testSubscriber) OnNamespaceCreated(namespace Namespace) {
	s.namespaces = append(s.namespaces, namespace.Key())
}

func (s *testSubscriber) OnTypeCreated(t *Type) {
	s.types = append(s.types, t)
}
