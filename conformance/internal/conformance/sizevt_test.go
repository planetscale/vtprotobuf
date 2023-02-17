// Commercial secret, LLC "RevTech". Refer to CONFIDENTIAL file in the root for details

package conformance

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSizeVTBasic(t *testing.T) {
	msg := &TestAllTypesProto3{
		OptionalTimestamp: timestamppb.Now(),
	}
	require.Equal(t, proto.Size(msg), msg.SizeVT())
}
