package editions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestOpaqueRoundTrip(t *testing.T) {
	for _, msg := range []*EditionsMessage{
		EditionsMessage_builder{RequiredInt32: proto.Int32(0)}.Build(),
		EditionsMessage_builder{
			ExplicitInt32:  proto.Int32(1),
			ExplicitString: proto.String("explicit"),
			ImplicitInt32:  2,
			ImplicitString: "implicit",
			RequiredInt32:  proto.Int32(3),
		}.Build(),
		// Explicit presence of a zero value must survive the round trip.
		EditionsMessage_builder{
			ExplicitInt32:  proto.Int32(0),
			ExplicitString: proto.String(""),
			RequiredInt32:  proto.Int32(0),
		}.Build(),
	} {
		vt, err := msg.MarshalVT()
		require.NoError(t, err)
		require.Equal(t, msg.SizeVT(), len(vt))

		var got EditionsMessage
		require.NoError(t, got.UnmarshalVT(vt))
		require.True(t, proto.Equal(msg, &got))
		require.True(t, msg.EqualVT(msg.CloneVT()))
	}
}

func TestOpenRoundTrip(t *testing.T) {
	for _, msg := range []*OpenMessage{
		{RequiredInt32: proto.Int32(0)},
		{
			ExplicitInt32:  proto.Int32(1),
			ExplicitString: proto.String("explicit"),
			ImplicitInt32:  2,
			ImplicitString: "implicit",
			RequiredInt32:  proto.Int32(3),
		},
		{
			ExplicitInt32:  proto.Int32(0),
			ExplicitString: proto.String(""),
			RequiredInt32:  proto.Int32(0),
		},
	} {
		vt, err := msg.MarshalVT()
		require.NoError(t, err)
		require.Equal(t, msg.SizeVT(), len(vt))

		var got OpenMessage
		require.NoError(t, got.UnmarshalVT(vt))
		require.True(t, proto.Equal(msg, &got))
		require.True(t, msg.EqualVT(msg.CloneVT()))
	}
}

func TestImplicitPresence(t *testing.T) {
	require.True(t, (&OpenMessage{ImplicitInt32: 0, RequiredInt32: proto.Int32(0)}).
		EqualVT(&OpenMessage{RequiredInt32: proto.Int32(0)}))
	require.False(t, (&OpenMessage{ExplicitInt32: proto.Int32(0), RequiredInt32: proto.Int32(0)}).
		EqualVT(&OpenMessage{RequiredInt32: proto.Int32(0)}))
}

func TestRequiredFieldUnset(t *testing.T) {
	_, err := (&OpenMessage{}).MarshalVT()
	require.ErrorContains(t, err, "required field required_int32 not set")

	_, err = (&EditionsMessage{}).MarshalVT()
	require.ErrorContains(t, err, "required field required_int32 not set")
}
