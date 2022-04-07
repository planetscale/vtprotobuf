package pool

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
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
