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
		Foo1: "⚡速いA6@wdkyAiX!7Ls7Tp9_NÉŠŻ⚡",
		Foo2: []byte{8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5, 6, 7, 8},
	}

	t0Bytes, err := t0.MarshalVT()
	require.NoError(t, err)

	// Note: this test performs checks on the Data uintptr of headers. It works as long as the garbage collector doesn't
	// move memory. To provide guarantee that this test works, consider using https://pkg.go.dev/runtime#Pinner when
	// upgrading Go to >= 1.21.
	t0BytesStart := (*reflect.SliceHeader)(unsafe.Pointer(&t0Bytes)).Data
	t0BytesEnd := t0BytesStart + uintptr(t0.SizeVT()) - 1

	// Helper functions to check that the underlying array belongs to t0Bytes
	assertStringBelongs := func(s string, belongs bool) {
		hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
		start := hdr.Data
		end := start + uintptr(len(s)) - 1
		assert.Equal(t, belongs, start >= t0BytesStart && start < t0BytesEnd)
		assert.Equal(t, belongs, end > t0BytesStart && end <= t0BytesEnd)
	}
	assertBytesBelongs := func(s []byte, belongs bool) {
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		start := hdr.Data
		end := start + uintptr(len(s)) - 1
		assert.Equal(t, belongs, start >= t0BytesStart && start < t0BytesEnd)
		assert.Equal(t, belongs, end > t0BytesStart && end <= t0BytesEnd)
	}

	t1Safe := &UnsafeTest{}
	require.NoError(t, t1Safe.UnmarshalVT(t0Bytes))
	require.NotEqual(t, uintptr(unsafe.Pointer(t0)), uintptr(unsafe.Pointer(t1Safe)))

	assertStringBelongs(t1Safe.Foo1, false)
	assertBytesBelongs(t1Safe.Foo2, false)

	t1Unsafe := &UnsafeTest{}
	require.NoError(t, t1Unsafe.UnmarshalVTUnsafe(t0Bytes))
	require.NotEqual(t, uintptr(unsafe.Pointer(t0)), uintptr(unsafe.Pointer(t1Unsafe)))

	assertStringBelongs(t1Unsafe.Foo1, true)
	assertBytesBelongs(t1Unsafe.Foo2, true)
}
