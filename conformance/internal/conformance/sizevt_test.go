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

func TestSizeVTBasic(t *testing.T) {
	now := time.Now()
	msg := &TestAllTypesProto3{
		OptionalTimestamp: timestamppb.New(now),
		OptionalDuration:  durationpb.New(time.Since(now)),
	}
	require.Equal(t, proto.Size(msg), msg.SizeVT())
}