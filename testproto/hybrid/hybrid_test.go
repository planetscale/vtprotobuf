//go:build !protoopaque

package hybrid

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestRoundTrip(t *testing.T) {
	msg := &HybridMessage{
		ImplicitInt32:  1,
		OptionalString: proto.String("s"),
		Children:       []*HybridMessage{{ImplicitInt32: 2}},
	}

	vt, err := msg.MarshalVT()
	require.NoError(t, err)
	require.Equal(t, msg.SizeVT(), len(vt))

	var got HybridMessage
	require.NoError(t, got.UnmarshalVT(vt))
	require.True(t, proto.Equal(msg, &got))
	require.True(t, msg.EqualVT(msg.CloneVT()))
}
