package unique

import (
	"maps"
	"slices"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalSameMemory(t *testing.T) {
	m := &UniqueFieldExtension{
		Foo: "bar",
		Bar: map[string]int64{"key": 100},
		Baz: map[int64]string{100: "value"},
	}

	b, err := m.MarshalVTStrict()
	require.NoError(t, err)

	m2 := &UniqueFieldExtension{}
	require.NoError(t, m2.UnmarshalVT(b))

	m3 := &UniqueFieldExtension{}
	require.NoError(t, m3.UnmarshalVT(b))

	require.Equal(t, unsafe.StringData(m2.Foo), unsafe.StringData(m3.Foo))

	keys2 := slices.Collect(maps.Keys(m2.Bar))
	keys3 := slices.Collect(maps.Keys(m3.Bar))
	require.Len(t, keys2, 1)
	require.Len(t, keys3, 1)
	require.Equal(t, unsafe.StringData(keys2[0]), unsafe.StringData(keys3[0]))

	values2 := slices.Collect(maps.Values(m2.Baz))
	values3 := slices.Collect(maps.Values(m3.Baz))
	require.Len(t, values2, 1)
	require.Len(t, values2, 1)
	require.Equal(t, unsafe.StringData(values2[0]), unsafe.StringData(values3[0]))
}
