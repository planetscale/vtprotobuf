package conformance

import (
	"github.com/planetscale/vtprotobuf/types/known/structpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	upstreamstructpb "google.golang.org/protobuf/types/known/structpb"
	"testing"
)

func TestEmptyOneof(t *testing.T) {
	// Regression test for https://github.com/planetscale/vtprotobuf/issues/61
	t.Run("all proto", func(t *testing.T) {
		msg := &TestAllTypesProto3{OneofField: &TestAllTypesProto3_OneofNestedMessage{}}
		upstream, _ := proto.Marshal(msg)
		vt, _ := msg.MarshalVTStrict()
		require.Equal(t, upstream, vt)
	})
	t.Run("structpb", func(t *testing.T) {
		msg := &structpb.Value{Kind: &upstreamstructpb.Value_StructValue{}}
		upstream, _ := proto.Marshal((*upstreamstructpb.Value)(msg))
		vt, _ := msg.MarshalVTStrict()
		require.Equal(t, upstream, vt)
	})
}
