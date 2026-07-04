package grpc

import (
	"fmt"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

type Codec struct{}

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

type reseter interface{
	Reset()
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
	//All types that implement github.com/golang/protobuf/proto.Message have a Reset method
	vv, ok := v.(reseter)
	if !ok {
		return fmt.Errorf("failed to unmarshal: can't reset. Message type %T don't implement Reset()", vv)
	}
	vv.Reset()
	return vt.UnmarshalVT(data)
}

func (Codec) Name() string {
	return Name
}
