package unsafe

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UnmarshalVTUnsafe(t *testing.T) {
	t0 := &UnsafeTest{
		// TODO: add more test cases, consider fuzzing
		Foo1: "c4c5cVG%$VG$%_KCREÉŠŻ",
		Foo2: []byte{1, 2, 3, 4, 5, 6, 7, 8},
	}

	t0Bytes, err := t0.MarshalVT()
	require.NoError(t, err)

	t1 := &UnsafeTest{}
	require.NoError(t, t1.UnmarshalVTUnsafe(t0Bytes))
	require.NotEqual(t, uintptr(unsafe.Pointer(t0)), uintptr(unsafe.Pointer(t1)))

	t0BytesStart := (*reflect.SliceHeader)(unsafe.Pointer(&t0Bytes)).Data
	t0BytesEnd := t0BytesStart + uintptr(t0.SizeVT()) - 1

	assert.Equal(t, t0.Foo1, t1.Foo1)
	hdr0s := (*reflect.StringHeader)(unsafe.Pointer(&t0.Foo1))
	hdr1s := (*reflect.StringHeader)(unsafe.Pointer(&t1.Foo1))
	assert.False(t, hdr0s.Data > t0BytesStart && hdr0s.Data < t0BytesEnd)
	// TODO: maybe add something for UnmarshalVT
	// the underlying data of Foo1 belongs to t0Bytes because it is unsafe
	assert.True(t, hdr1s.Data > t0BytesStart && hdr1s.Data < t0BytesEnd)

	assert.Equal(t, t0.Foo2, t1.Foo2)
	hdr0b := (*reflect.SliceHeader)(unsafe.Pointer(&t0.Foo2))
	hdr1b := (*reflect.SliceHeader)(unsafe.Pointer(&t1.Foo2))
	assert.False(t, hdr0b.Data > t0BytesStart && hdr0b.Data < t0BytesEnd)
	// the underlying data of Foo2 belongs to t0Bytes because it is unsafe
	// TODO: it fails here, there is an issue with bytes
	assert.True(t, hdr1b.Data > t0BytesStart && hdr1b.Data < t0BytesEnd)
}
