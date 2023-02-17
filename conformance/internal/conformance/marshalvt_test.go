// Commercial secret, LLC "RevTech". Refer to CONFIDENTIAL file in the root for details

package conformance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMarshalVTBasic(t *testing.T) {
	msg := &TestAllTypesProto3{
		OptionalTimestamp: timestamppb.Now(),
		OptionalDuration:  durationpb.New(time.Second),
	}
	serializedOrig, err := proto.Marshal(msg)
	require.NoError(t, err)
	serializedVT, err := msg.MarshalVT()
	require.NoError(t, err)
	require.Equal(t, serializedOrig, serializedVT)
}
