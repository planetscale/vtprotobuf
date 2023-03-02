// Commercial secret, LLC "RevTech". Refer to CONFIDENTIAL file in the root for details

package conformance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUnmarshalVTBasic(t *testing.T) {
	msg := &TestAllTypesProto3{
		OptionalTimestamp: timestamppb.Now(),
		OptionalDuration:  durationpb.New(time.Second),
		OptionalValue:     structpb.NewStringValue("kek"),
	}
	serializedOrig, err := proto.Marshal(msg)
	require.NoError(t, err)

	got := &TestAllTypesProto3{}
	require.NoError(t, got.UnmarshalVT(serializedOrig))
	require.True(t, proto.Equal(msg, got))
	require.True(t, msg.EqualVT(got))
}
