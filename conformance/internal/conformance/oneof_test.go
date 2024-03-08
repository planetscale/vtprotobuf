package conformance

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestEmptyOneoff(t *testing.T) {
	// Regression test for https://github.com/planetscale/vtprotobuf/issues/61
	msg := &TestAllTypesProto3{OneofField: &TestAllTypesProto3_OneofNestedMessage{}}
	upstream, _ := proto.Marshal(msg)
	vt, _ := msg.MarshalVTStrict()
	require.Equal(t, upstream, vt)
}
