package pool

import (
	"google.golang.org/protobuf/proto"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Pool_slice_data_override(t *testing.T) {
	dataRequest := Test1{
		Sl: []string{"a", "b", "c"},
	}
	dataReqBytes, err := dataRequest.MarshalVT()
	require.NoError(t, err)

	nilRequest := Test1{
		Sl: []string{},
	}
	nilReqBytes, err := nilRequest.MarshalVT()
	require.NoError(t, err)

	req := Test1FromVTPool()
	err = req.UnmarshalVT(dataReqBytes)
	require.NoError(t, err)
	req.ReturnToVTPool()

	req = Test1FromVTPool()
	err = req.UnmarshalVT(nilReqBytes)
	require.NoError(t, err)

	assert.Equal(t, []string{}, req.Sl)
}

func Test_Pool_slice_recreation(t *testing.T) {
	var b int32 = 10
	dataRequest := Test2{
		Sl: []*Slice2{
			{
				A: map[int64]int64{
					22: 33,
				},
				B: &b,
				C: []string{"one", "two"},
				D: &Element2{
					A: 10,
				},
				E: "something",
				F: 123,
			},
		},
	}
	dataReqBytes, err := dataRequest.MarshalVT()
	require.NoError(t, err)

	nilRequest := Test2{
		Sl: []*Slice2{
			{
				A: nil,
				B: nil,
				C: nil,
				D: nil,
				E: "",
				F: 0,
			},
		},
	}
	nilReqBytes, err := nilRequest.MarshalVT()
	require.NoError(t, err)

	req := Test2FromVTPool()
	err = req.UnmarshalVT(dataReqBytes)
	log.Println(req)
	require.NoError(t, err)
	req.ReturnToVTPool()

	req = Test2FromVTPool()
	err = req.UnmarshalVT(nilReqBytes)
	require.NoError(t, err)

	assert.Nil(t, req.Sl[0].A)
	assert.Nil(t, req.Sl[0].B)
	assert.Nil(t, req.Sl[0].D)
	assert.Nil(t, req.Sl[0].C)
	assert.Zero(t, req.Sl[0].E)
	assert.Zero(t, req.Sl[0].F)
}

func Test_Pool_Oneof(t *testing.T) {
	t1_t_t1 := &OneofTest_Test1{
		A: 42,
	}
	t1_t := &OneofTest_Test1_{
		Test1: t1_t_t1,
	}
	t1 := &OneofTest{
		Test: t1_t,
	}

	t2_t_t2 := &OneofTest_Test2{
		B: []string{"str1", "str2", "str3"},
	}
	t2_t := &OneofTest_Test2_{
		Test2: t2_t_t2,
	}
	t2 := &OneofTest{
		Test: t2_t,
	}

	t3_t_t3_c := &OneofTest_Test3_Element2{
		D: 51,
	}
	t3_t_t3 := &OneofTest_Test3{
		C: t3_t_t3_c,
	}
	t3_t := &OneofTest_Test3_{
		Test3: t3_t_t3,
	}
	t3 := &OneofTest{
		Test: t3_t,
	}

	t4 := &OneofTest{Test: nil}

	t1Bytes, err := t1.MarshalVT()
	require.NoError(t, err)
	t2Bytes, err := t2.MarshalVT()
	require.NoError(t, err)
	t3Bytes, err := t3.MarshalVT()
	require.NoError(t, err)
	t4Bytes, err := t4.MarshalVT()
	require.NoError(t, err)

	t1.ReturnToVTPool()

	t5 := OneofTestFromVTPool()
	require.NoError(t, t5.UnmarshalVT(t1Bytes))
	require.Equal(t, &t1, &t5)
	require.Equal(t, t1_t, t5.Test.(*OneofTest_Test1_))
	require.Equal(t, t1_t_t1, t5.Test.(*OneofTest_Test1_).Test1)
	require.Equal(t, &t1_t_t1, &t5.Test.(*OneofTest_Test1_).Test1)

	t2.ReturnToVTPool()
	t6 := OneofTestFromVTPool()
	require.NoError(t, t6.UnmarshalVT(t2Bytes))
	require.Equal(t, &t2, &t6)
	require.Equal(t, t2_t, t6.Test.(*OneofTest_Test2_))
	require.Equal(t, t2_t_t2, t6.Test.(*OneofTest_Test2_).Test2)
	require.Equal(t, &t2_t_t2, &t6.Test.(*OneofTest_Test2_).Test2)

	t3.ReturnToVTPool()
	t7 := OneofTestFromVTPool()
	require.NoError(t, t7.UnmarshalVT(t3Bytes))
	require.Equal(t, &t3, &t7)
	require.Equal(t, t3_t, t7.Test.(*OneofTest_Test3_))
	require.Equal(t, t3_t_t3, t7.Test.(*OneofTest_Test3_).Test3)
	require.Equal(t, &t3_t_t3, &t7.Test.(*OneofTest_Test3_).Test3)
	require.Equal(t, &t3_t_t3_c, &t7.Test.(*OneofTest_Test3_).Test3.C)

	t4.ReturnToVTPool()
	t8 := OneofTestFromVTPool()
	require.NoError(t, t8.UnmarshalVT(t4Bytes))
	require.Equal(t, &t4, &t8)
}

func Test_Pool_Optional(t *testing.T) {
	m := &MemoryPoolExtension{
		Foo1: "foo1",
		Foo2: 123,
		Foo3: &OptionalMessage{},
	}

	mBytes, err := m.MarshalVT()
	require.NoError(t, err)

	mUnmarshal := &MemoryPoolExtension{}
	err = proto.Unmarshal(mBytes, mUnmarshal)
	require.NoError(t, err)

	require.True(t, m.EqualVT(mUnmarshal))

	m.ReturnToVTPool()
	mFromPool := MemoryPoolExtensionFromVTPool()
	require.True(t, mFromPool.EqualVT(&MemoryPoolExtension{}))
}
