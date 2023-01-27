package proto3opt

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestEmptyBytesMarshalling(t *testing.T) {
	a := &OptionalFieldInProto3{
		OptionalBytes: nil,
	}
	b := &OptionalFieldInProto3{
		OptionalBytes: []byte{},
	}

	type Message interface {
		proto.Message
		MarshalVT() ([]byte, error)
	}

	for _, msg := range []Message{a, b} {
		vt, err := msg.MarshalVT()
		require.NoError(t, err)
		goog, err := proto.Marshal(msg)
		require.NoError(t, err)
		require.Equal(t, goog, vt)
	}
}
