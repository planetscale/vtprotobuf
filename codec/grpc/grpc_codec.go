package grpc

import "fmt"

// Name is the name registered for the proto compressor.
const Name = "proto"

type Codec struct{}

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

// vtprotoSlabMessage is implemented by messages that opted into slab
// unmarshalling (`option (vtproto.slab) = true;` with the unmarshal_slab
// feature): the codec prefers the slab entry point, so opted-in messages get
// arena-backed decoding without any call-site changes.
type vtprotoSlabMessage interface {
	UnmarshalVTSlab([]byte) error
}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	vt, ok := v.(vtprotoMessage)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T (missing vtprotobuf helpers)", v)
	}
	return vt.MarshalVT()
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	if vt, ok := v.(vtprotoSlabMessage); ok {
		return vt.UnmarshalVTSlab(data)
	}
	vt, ok := v.(vtprotoMessage)
	if !ok {
		return fmt.Errorf("failed to unmarshal, message is %T (missing vtprotobuf helpers)", v)
	}
	return vt.UnmarshalVT(data)
}

func (Codec) Name() string {
	return Name
}
