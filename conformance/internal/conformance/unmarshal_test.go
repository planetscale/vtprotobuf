package conformance

import (
	"testing"
	utf8 "unicode/utf8"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestUnmarshalVTInvalidUTF8(t *testing.T) {
	stringPtr := func(x string) *string { return &x }
	payload := []byte{0x12, 0x34, 0x56, 0x78, 0x9a}
	require.False(t, utf8.Valid(payload))
	{
		msg := &TestAllTypesProto2{
			OptionalString: stringPtr(string(payload)),
		}
		data, err := proto.Marshal(msg)
		require.NoError(t, err)

		var msg2 TestAllTypesProto2
		err = msg2.UnmarshalVT(data)
		require.Equal(t, "proto: invalid UTF-8 string", err.Error())
	}
	{
		msg := &TestAllTypesProto2{
			MapStringString: map[string]string{
				"key": string(payload),
			},
		}
		data, err := proto.Marshal(msg)
		require.NoError(t, err)

		var msg2 TestAllTypesProto2
		err = msg2.UnmarshalVT(data)
		require.Equal(t, "proto: invalid UTF-8 string", err.Error())
	}
}
