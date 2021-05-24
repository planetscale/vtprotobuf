package drpc

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

func Marshal(msg interface{}) ([]byte, error) {
	return msg.(vtprotoMessage).MarshalVT()
}

func Unmarshal(buf []byte, msg interface{}) error {
	return msg.(vtprotoMessage).UnmarshalVT(buf)
}

func JSONMarshal(msg interface{}) ([]byte, error) {
	return protojson.Marshal(msg.(proto.Message))
}

func JSONUnmarshal(buf []byte, msg interface{}) error {
	return protojson.Unmarshal(buf, msg.(proto.Message))
}
