package wkt

import (
	"google.golang.org/protobuf/proto"
	"testing"
	"time"

	any "github.com/planetscale/vtprotobuf/types/known/any"
	duration "github.com/planetscale/vtprotobuf/types/known/duration"
	empty "github.com/planetscale/vtprotobuf/types/known/empty"
	field_mask "github.com/planetscale/vtprotobuf/types/known/field_mask"
	timestamp "github.com/planetscale/vtprotobuf/types/known/timestamp"
	wrappers "github.com/planetscale/vtprotobuf/types/known/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWellKnownTypes(t *testing.T) {
	dur := duration.New(4*time.Hour + 2*time.Second)

	anyVal, err := any.New(dur)
	require.NoError(t, err)

	fieldMask, err := field_mask.New(dur, "seconds")
	require.NoError(t, err)

	m := &MessageWithWKT{
		Any:         anyVal,
		Duration:    dur,
		Empty:       &empty.Empty{},
		FieldMask:   fieldMask,
		Timestamp:   timestamp.Now(),
		DoubleValue: wrappers.Double(123456789.123456789),
		FloatValue:  wrappers.Float(123456789.123456789),
		Int64Value:  wrappers.Int64(123456789),
		Uint64Value: wrappers.UInt64(123456789),
		Int32Value:  wrappers.Int32(123456789),
		Uint32Value: wrappers.UInt32(123456789),
		BoolValue:   wrappers.Bool(true),
		StringValue: wrappers.String("String marshalling and unmarshalling test"),
		BytesValue:  wrappers.Bytes([]byte("Bytes marshalling and unmarshalling test")),
	}

	golangBytes, err := proto.Marshal(m)
	require.NoError(t, err)

	vtProtoBytes, err := m.MarshalVT()
	require.NoError(t, err)

	require.NotEmpty(t, golangBytes)
	require.NotEmpty(t, vtProtoBytes)
	assert.Equal(t, golangBytes, vtProtoBytes)

	var (
		golangMsg  = &MessageWithWKT{}
		vtProtoMsg = &MessageWithWKT{}
	)

	require.NoError(t, proto.Unmarshal(golangBytes, golangMsg))
	require.NoError(t, vtProtoMsg.UnmarshalVT(vtProtoBytes))

	assert.Equal(t, golangMsg.String(), vtProtoMsg.String())
}
