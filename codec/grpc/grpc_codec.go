package grpc

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

type Codec struct{}

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	vt, ok := v.(vtprotoMessage)
	if ok {
		return vt.MarshalVT()
	}

	// fallback to native Protobuf format
	vv, ok := v.(proto.Message)
	if ok {
		return proto.Marshal(vv)
	}

	// return error if neither can marshal
	return nil, fmt.Errorf("failed to marshal, message is %T, tried vtproto and proto", v)
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	vt, ok := v.(vtprotoMessage)
	if ok {
		return vt.UnmarshalVT(data)
	}

	// fallback to native Protobuf format
	vv, ok := v.(proto.Message)
	if ok {
		return proto.Unmarshal(data, vv)
	}

	// return error if neither can unmarshal
	return fmt.Errorf("failed to unmarshal, message is %T, tried vtproto and proto", v)
}

func (Codec) Name() string {
	return Name
}
