package opaque

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func messages() []*OpaqueMessage {
	return []*OpaqueMessage{
		{},
		OpaqueMessage_builder{
			ImplicitInt32:  1,
			ImplicitString: "implicit",
			OptionalInt32:  proto.Int32(-2),
			OptionalDouble: proto.Float64(3.5),
			OptionalBool:   proto.Bool(true),
			OptionalEnum:   OpaqueEnum_OPAQUE_ENUM_ONE.Enum(),
			OptionalString: proto.String("optional"),
			OptionalBytes:  []byte("optional bytes"),
			Inner:          Inner_builder{Value: 4}.Build(),
			RepeatedInt32:  []int32{5, 6},
			RepeatedInner: []*Inner{
				Inner_builder{Value: 7}.Build(),
				Inner_builder{Value: 8}.Build(),
			},
			MapInt32:   map[string]int32{"k": 9},
			OneofInt32: proto.Int32(10),
		}.Build(),
		// Only the presence bits distinguish these from an empty message.
		OpaqueMessage_builder{
			OptionalInt32:  proto.Int32(0),
			OptionalBool:   proto.Bool(false),
			OptionalString: proto.String(""),
			OptionalBytes:  []byte{},
		}.Build(),
		OpaqueMessage_builder{OneofString: proto.String("s")}.Build(),
		OpaqueMessage_builder{OneofInner: Inner_builder{Value: 11}.Build()}.Build(),
	}
}

func TestRoundTrip(t *testing.T) {
	for _, msg := range messages() {
		vt, err := msg.MarshalVT()
		require.NoError(t, err)
		require.Equal(t, msg.SizeVT(), len(vt))

		var fromVT OpaqueMessage
		require.NoError(t, proto.Unmarshal(vt, &fromVT))
		require.True(t, proto.Equal(msg, &fromVT))

		enc, err := proto.Marshal(msg)
		require.NoError(t, err)

		var toVT OpaqueMessage
		require.NoError(t, toVT.UnmarshalVT(enc))
		require.True(t, proto.Equal(msg, &toVT))

		require.True(t, msg.EqualVT(msg.CloneVT()))
		require.True(t, proto.Equal(msg, msg.CloneVT()))
	}
}

func TestPresenceOnlyDiffers(t *testing.T) {
	set := OpaqueMessage_builder{OptionalInt32: proto.Int32(0)}.Build()
	require.False(t, set.EqualVT(&OpaqueMessage{}))
	require.True(t, set.EqualVT(OpaqueMessage_builder{OptionalInt32: proto.Int32(0)}.Build()))
}

func TestPool(t *testing.T) {
	msg := PooledOpaque_builder{
		Payload:       []byte("payload"),
		Counter:       proto.Int32(1),
		RepeatedInner: []*Inner{Inner_builder{Value: 2}.Build()},
	}.Build()
	enc, err := proto.Marshal(msg)
	require.NoError(t, err)

	pooled := PooledOpaqueFromVTPool()
	require.NoError(t, pooled.UnmarshalVT(enc))
	require.True(t, proto.Equal(msg, pooled))
	pooled.ReturnToVTPool()

	// A recycled message must not carry over presence bits the next payload omits.
	reused := PooledOpaqueFromVTPool()
	require.NoError(t, reused.UnmarshalVT(nil))
	require.False(t, reused.HasCounter())
	require.True(t, proto.Equal(&PooledOpaque{}, reused))
	reused.ReturnToVTPool()
}
