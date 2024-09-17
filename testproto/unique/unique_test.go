package unique

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalSameMemory(t *testing.T) {
	m := &UniqueFieldExtension{
		Foo: "bar",
	}

	b, err := m.MarshalVTStrict()
	require.NoError(t, err)

	m2 := &UniqueFieldExtension{}
	require.NoError(t, m2.UnmarshalVT(b))

	m3 := &UniqueFieldExtension{}
	require.NoError(t, m3.UnmarshalVT(b))

	require.Equal(t, unsafe.StringData(m2.Foo), unsafe.StringData(m3.Foo))
}
