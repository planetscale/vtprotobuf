package wkt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestWellKnownTypes(t *testing.T) {
	dur := durationpb.New(4*time.Hour + 2*time.Second)

	anyVal, err := anypb.New(dur)
	require.NoError(t, err)

	fieldMask, err := fieldmaskpb.New(dur, "seconds")
	require.NoError(t, err)

	m := &MessageWithWKT{
		Any:         anyVal,
		Duration:    dur,
		Empty:       &emptypb.Empty{},
		FieldMask:   fieldMask,
		Timestamp:   timestamppb.Now(),
		DoubleValue: wrapperspb.Double(123456789.123456789),
		FloatValue:  wrapperspb.Float(123456789.123456789),
		Int64Value:  wrapperspb.Int64(123456789),
		Uint64Value: wrapperspb.UInt64(123456789),
		Int32Value:  wrapperspb.Int32(123456789),
		Uint32Value: wrapperspb.UInt32(123456789),
		BoolValue:   wrapperspb.Bool(true),
		StringValue: wrapperspb.String("String marshalling and unmarshalling test"),
		BytesValue:  wrapperspb.Bytes([]byte("Bytes marshalling and unmarshalling test")),
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
