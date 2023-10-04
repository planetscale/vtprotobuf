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
	originalData, err := t0.MarshalVT()
	require.NoError(t, err)

	// Note: this test performs checks on the Data uintptr of headers. It works as long as the garbage collector doesn't
	// move memory. To provide guarantee that this test works, consider using https://pkg.go.dev/runtime#Pinner when
	// upgrading Go to >= 1.21.
	originalStart := (*reflect.SliceHeader)(unsafe.Pointer(&originalData)).Data
	originalEnd := originalStart + uintptr(len(originalData)) - 1

	// Helper functions to check whether the underlying array belongs to originalData
	assertStringIsOriginal := func(s string, isOriginal bool) {
		hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
		start := hdr.Data
		end := start + uintptr(len(s)) - 1
		assert.Equal(t, isOriginal, start >= originalStart && start < originalEnd)
		assert.Equal(t, isOriginal, end > originalStart && end <= originalEnd)
	}
	assertBytesIsOriginal := func(s []byte, isOriginal bool) {
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		start := hdr.Data
		end := start + uintptr(len(s)) - 1
		assert.Equal(t, isOriginal, start >= originalStart && start < originalEnd)
		assert.Equal(t, isOriginal, end > originalStart && end <= originalEnd)
	}

	t1Safe := &UnsafeTest{}
	require.NoError(t, t1Safe.UnmarshalVT(originalData))

	assertStringIsOriginal(t1Safe.Foo1, false)
	assertBytesIsOriginal(t1Safe.Foo2, false)

	t1Unsafe := &UnsafeTest{}
	require.NoError(t, t1Unsafe.UnmarshalVTUnsafe(originalData))

	assertStringIsOriginal(t1Unsafe.Foo1, true)
	assertBytesIsOriginal(t1Unsafe.Foo2, true)
}
