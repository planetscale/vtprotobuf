package pool

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func Test_Pool(t *testing.T) {

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

}
