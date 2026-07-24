package drpc

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

// vtprotoSlabMessage is implemented by messages that opted into slab
// unmarshalling (`option (vtproto.slab) = true;` with the unmarshal_slab
// feature); the codec prefers the slab entry point when it is available.
type vtprotoSlabMessage interface {
	UnmarshalVTSlab([]byte) error
}

func Marshal(msg interface{}) ([]byte, error) {
	return msg.(vtprotoMessage).MarshalVT()
}

func Unmarshal(buf []byte, msg interface{}) error {
	if m, ok := msg.(vtprotoSlabMessage); ok {
		return m.UnmarshalVTSlab(buf)
	}
	return msg.(vtprotoMessage).UnmarshalVT(buf)
}

func JSONMarshal(msg interface{}) ([]byte, error) {
	return protojson.Marshal(msg.(proto.Message))
}

func JSONUnmarshal(buf []byte, msg interface{}) error {
	return protojson.Unmarshal(buf, msg.(proto.Message))
}
