package unsafe

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertStringIsOriginal asserts whether the underlying array of s belongs to originalData.
//
// Note: performing checks on the uintptr to the underlying arrays works as long as the garbage
// collector doesn't move memory. To provide guarantee that this test works, consider using
// https://pkg.go.dev/runtime#Pinner when upgrading Go to >= 1.21.
func assertStringIsOriginal(t *testing.T, s string, belongs bool, originalData []byte) {
	start := uintptr(unsafe.Pointer(unsafe.StringData(s)))
	// empty string has no underlying array, compare pointer to nil
	if len(s) == 0 {
		assert.Equal(t, uintptr(unsafe.Pointer(nil)), start)
		return
	}
	end := start + uintptr(len(s)) - 1

	originalStart := uintptr(unsafe.Pointer(unsafe.SliceData(originalData)))
	originalEnd := originalStart + uintptr(len(originalData)) - 1

	assert.Equal(t, belongs, start >= originalStart && start < originalEnd)
	assert.Equal(t, belongs, end > originalStart && end <= originalEnd)
}

// assertBytesAreOriginal is the same as assertStringIsOriginal for a []byte.
func assertBytesAreOriginal(t *testing.T, b []byte, belongs bool, originalData []byte) {
	originalStart := uintptr(unsafe.Pointer(unsafe.SliceData(originalData)))
	originalEnd := originalStart + uintptr(len(originalData)) - 1

	start := uintptr(unsafe.Pointer(unsafe.SliceData(b)))
	end := start + uintptr(len(b)) - 1
	assert.Equal(t, belongs, start >= originalStart && start < originalEnd)
	assert.Equal(t, belongs, end > originalStart && end <= originalEnd)
}

func Test_UnmarshalVTUnsafe(t *testing.T) {
	testString := "⚡速いA6@wdkyAiX!7Ls7Tp9_NÉŠŻ⚡"
	testBytes := []byte{8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5, 6, 7, 8}

	// simple case
	t1Orig := &UnsafeTest{
		Sub: &UnsafeTest_Sub1_{
			Sub1: &UnsafeTest_Sub1{
				S: testString,
				B: testBytes,
			},
		},
	}
	originalData1, err := t1Orig.MarshalVT()
	require.NoError(t, err)

	t1Safe := &UnsafeTest{}
	require.NoError(t, t1Safe.UnmarshalVT(originalData1))
	assert.Equal(t, (t1Orig.Sub).(*UnsafeTest_Sub1_).Sub1.S, (t1Safe.Sub).(*UnsafeTest_Sub1_).Sub1.S)
	assert.Equal(t, (t1Orig.Sub).(*UnsafeTest_Sub1_).Sub1.B, (t1Safe.Sub).(*UnsafeTest_Sub1_).Sub1.B)
	assertStringIsOriginal(t, (t1Safe.Sub).(*UnsafeTest_Sub1_).Sub1.S, false, originalData1)
	assertBytesAreOriginal(t, (t1Safe.Sub).(*UnsafeTest_Sub1_).Sub1.B, false, originalData1)

	t1Unsafe := &UnsafeTest{}
	require.NoError(t, t1Unsafe.UnmarshalVTUnsafe(originalData1))
	assert.Equal(t, (t1Orig.Sub).(*UnsafeTest_Sub1_).Sub1.S, (t1Unsafe.Sub).(*UnsafeTest_Sub1_).Sub1.S)
	assert.Equal(t, (t1Orig.Sub).(*UnsafeTest_Sub1_).Sub1.B, (t1Unsafe.Sub).(*UnsafeTest_Sub1_).Sub1.B)
	assertStringIsOriginal(t, (t1Unsafe.Sub).(*UnsafeTest_Sub1_).Sub1.S, true, originalData1)
	assertBytesAreOriginal(t, (t1Unsafe.Sub).(*UnsafeTest_Sub1_).Sub1.B, true, originalData1)

	// repeated field
	t2Orig := &UnsafeTest{
		Sub: &UnsafeTest_Sub2_{
			Sub2: &UnsafeTest_Sub2{
				S: []string{testString, testString, testString},
				B: [][]byte{testBytes, testBytes, testBytes},
			},
		},
	}
	originalData2, err := t2Orig.MarshalVT()
	require.NoError(t, err)

	t2Safe := &UnsafeTest{}
	require.NoError(t, t2Safe.UnmarshalVT(originalData2))
	assert.Equal(t, (t2Orig.Sub).(*UnsafeTest_Sub2_).Sub2.S, (t2Safe.Sub).(*UnsafeTest_Sub2_).Sub2.S)
	assert.Equal(t, (t2Orig.Sub).(*UnsafeTest_Sub2_).Sub2.B, (t2Safe.Sub).(*UnsafeTest_Sub2_).Sub2.B)
	for i := 0; i < 3; i++ {
		assertStringIsOriginal(t, (t2Safe.Sub).(*UnsafeTest_Sub2_).Sub2.S[i], false, originalData2)
		assertBytesAreOriginal(t, (t2Safe.Sub).(*UnsafeTest_Sub2_).Sub2.B[i], false, originalData2)
	}

	t2Unsafe := &UnsafeTest{}
	require.NoError(t, t2Unsafe.UnmarshalVTUnsafe(originalData2))
	assert.Equal(t, (t2Orig.Sub).(*UnsafeTest_Sub2_).Sub2.S, (t2Unsafe.Sub).(*UnsafeTest_Sub2_).Sub2.S)
	assert.Equal(t, (t2Orig.Sub).(*UnsafeTest_Sub2_).Sub2.B, (t2Unsafe.Sub).(*UnsafeTest_Sub2_).Sub2.B)
	for i := 0; i < 3; i++ {
		assertStringIsOriginal(t, (t2Unsafe.Sub).(*UnsafeTest_Sub2_).Sub2.S[i], true, originalData2)
		assertBytesAreOriginal(t, (t2Unsafe.Sub).(*UnsafeTest_Sub2_).Sub2.B[i], true, originalData2)
	}

	// map[string]bytes field
	t3Orig := &UnsafeTest{
		Sub: &UnsafeTest_Sub3_{
			Sub3: &UnsafeTest_Sub3{
				Foo: map[string][]byte{testString: testBytes},
			},
		},
	}
	originalData3, err := t3Orig.MarshalVT()
	require.NoError(t, err)

	t3Safe := &UnsafeTest{}
	require.NoError(t, t3Safe.UnmarshalVT(originalData3))
	assert.Equal(t, (t3Orig.Sub).(*UnsafeTest_Sub3_).Sub3.Foo, (t3Safe.Sub).(*UnsafeTest_Sub3_).Sub3.Foo)
	for k, v := range (t3Safe.Sub).(*UnsafeTest_Sub3_).Sub3.Foo {
		assertStringIsOriginal(t, k, false, originalData3)
		assertBytesAreOriginal(t, v, false, originalData3)
	}

	t3Unsafe := &UnsafeTest{}
	require.NoError(t, t3Unsafe.UnmarshalVTUnsafe(originalData3))
	assert.Equal(t, (t3Orig.Sub).(*UnsafeTest_Sub3_).Sub3.Foo, (t3Unsafe.Sub).(*UnsafeTest_Sub3_).Sub3.Foo)
	for k, v := range (t3Unsafe.Sub).(*UnsafeTest_Sub3_).Sub3.Foo {
		assertStringIsOriginal(t, k, true, originalData3)
		assertBytesAreOriginal(t, v, true, originalData3)
	}

	// oneof field
	for _, stringVal := range []string{
		testString,
		"",
	} {
		t4OrigS := &UnsafeTest{
			Sub: &UnsafeTest_Sub4_{
				Sub4: &UnsafeTest_Sub4{
					Foo: &UnsafeTest_Sub4_S{
						S: stringVal,
					},
				},
			},
		}
		originalData4S, err := t4OrigS.MarshalVT()
		require.NoError(t, err)

		t4SafeS := &UnsafeTest{}
		require.NoError(t, t4SafeS.UnmarshalVT(originalData4S))
		assert.Equal(t, (t4OrigS.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_S).S, (t4SafeS.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_S).S)
		assertStringIsOriginal(t, (t4SafeS.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_S).S, false, originalData4S)

		t4UnsafeS := &UnsafeTest{}
		require.NoError(t, t4UnsafeS.UnmarshalVTUnsafe(originalData4S))
		assert.Equal(t, (t4OrigS.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_S).S, (t4UnsafeS.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_S).S)
		assertStringIsOriginal(t, (t4UnsafeS.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_S).S, true, originalData4S)
	}

	t4OrigB := &UnsafeTest{
		Sub: &UnsafeTest_Sub4_{
			Sub4: &UnsafeTest_Sub4{
				Foo: &UnsafeTest_Sub4_B{
					B: testBytes,
				},
			},
		},
	}
	originalData4B, err := t4OrigB.MarshalVT()
	require.NoError(t, err)

	t4SafeB := &UnsafeTest{}
	require.NoError(t, t4SafeB.UnmarshalVT(originalData4B))
	assert.Equal(t, (t4OrigB.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_B).B, (t4SafeB.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_B).B)
	assertBytesAreOriginal(t, (t4SafeB.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_B).B, false, originalData4B)

	t4UnsafeB := &UnsafeTest{}
	require.NoError(t, t4UnsafeB.UnmarshalVTUnsafe(originalData4B))
	assert.Equal(t, (t4OrigB.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_B).B, (t4UnsafeB.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_B).B)
	assertBytesAreOriginal(t, (t4UnsafeB.Sub).(*UnsafeTest_Sub4_).Sub4.Foo.(*UnsafeTest_Sub4_B).B, true, originalData4B)

	// map[string]string field
	for _, stringVal := range []string{
		testString,
		"",
	} {
		t5Orig := &UnsafeTest{
			Sub: &UnsafeTest_Sub5_{
				Sub5: &UnsafeTest_Sub5{
					Foo: map[string]string{testString: stringVal},
				},
			},
		}
		originalData5, err := t5Orig.MarshalVT()
		require.NoError(t, err)

		t5Safe := &UnsafeTest{}
		require.NoError(t, t5Safe.UnmarshalVT(originalData5))
		assert.Equal(t, (t5Orig.Sub).(*UnsafeTest_Sub5_).Sub5.Foo, (t5Safe.Sub).(*UnsafeTest_Sub5_).Sub5.Foo)
		for k, v := range (t5Safe.Sub).(*UnsafeTest_Sub5_).Sub5.Foo {
			assertStringIsOriginal(t, k, false, originalData5)
			assertStringIsOriginal(t, v, false, originalData5)
		}

		t5Unsafe := &UnsafeTest{}
		require.NoError(t, t5Unsafe.UnmarshalVTUnsafe(originalData5))
		assert.Equal(t, (t5Orig.Sub).(*UnsafeTest_Sub5_).Sub5.Foo, (t5Unsafe.Sub).(*UnsafeTest_Sub5_).Sub5.Foo)
		for k, v := range (t5Unsafe.Sub).(*UnsafeTest_Sub5_).Sub5.Foo {
			assertStringIsOriginal(t, k, true, originalData5)
			assertStringIsOriginal(t, v, true, originalData5)
		}
	}
}
