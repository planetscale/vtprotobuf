package grpc

import (
	"fmt"

	"github.com/golang/protobuf/proto" //nolint
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
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T (missing vtprotobuf helpers)", v)
	}
	return vt.MarshalVT()
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	vt, ok := v.(vtprotoMessage)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T (missing vtprotobuf helpers)", v)
	}
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T (can't reset)", vv)
	}
	vv.Reset()
	return vt.UnmarshalVT(data)
}

func (Codec) Name() string {
	return Name
}
