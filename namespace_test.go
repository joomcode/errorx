package errorx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	nsTest0             = NewNamespace("nsTest0")
	nsTest1             = NewNamespace("nsTest1")
	nsTest1Child        = nsTest1.NewSubNamespace("child")
	nsTestET0           = nsTest0.NewType("type0")
	nsTestET1           = nsTest1.NewType("type1")
	nsTestET1Child      = nsTestET1.NewSubtype("child")
	nsTestChild1ET      = nsTest1Child.NewType("type")
	nsTestChild1ETChild = nsTestChild1ET.NewSubtype("child")
)

func TestNamespaceName(t *testing.T) {
	require.EqualValues(t, "nsTest1", nsTest1.FullName())
	require.EqualValues(t, "nsTest1.child", nsTest1Child.FullName())
}

func TestNamespace(t *testing.T) {
	require.True(t, nsTest0.IsNamespaceOf(nsTestET0))
	require.False(t, nsTest1.IsNamespaceOf(nsTestET0))
	require.False(t, nsTest0.IsNamespaceOf(nsTestET1))
	require.True(t, nsTest1.IsNamespaceOf(nsTestET1))
}

func TestNamespaceSubtype(t *testing.T) {
	require.False(t, nsTest0.IsNamespaceOf(nsTestET1Child))
	require.True(t, nsTest1.IsNamespaceOf(nsTestET1Child))
}

func TestSubNamespace(t *testing.T) {
	require.False(t, nsTest1Child.IsNamespaceOf(nsTestET1))
	require.True(t, nsTest1Child.IsNamespaceOf(nsTestChild1ET))
	require.True(t, nsTest1Child.IsNamespaceOf(nsTestChild1ETChild))
}
