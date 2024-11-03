package conformance

import (
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"strings"
	"testing"
)

func roundTripUpstream(b []byte) ([]byte, error) {
	msg := &TestAllTypesProto3{}
	if err := proto.Unmarshal(b, msg); err != nil {
		return nil, err
	}
	res, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func roundTripVtprotobuf(b []byte) ([]byte, error) {
	msg := &TestAllTypesProto3{}
	if err := msg.UnmarshalVT(b); err != nil {
		return nil, err
	}
	res, err := msg.MarshalVT()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FuzzProto(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		u, uerr := roundTripUpstream(b)
		v, verr := roundTripVtprotobuf(b)
		if verr != nil && strings.Contains(verr.Error(), "wrong wireType") {
			t.Skip()
		}
		if uerr != nil && strings.Contains(uerr.Error(), "cannot parse invalid wire-format data") {
			t.Skip()
		}
		if (uerr != nil) != (verr != nil) {
			t.Fatalf("upstream err: %v (%v), vtprotobuf err: %v (%v)", uerr, u, verr, v)
		}
		us := &TestAllTypesProto3{}
		_ = proto.Unmarshal(b, us)
		us.unknownFields = nil

		t.Logf("upstream  : %v, %v", protojson.Format(us), prototext.Format(us))
		vt := &TestAllTypesProto3{}
		_ = vt.UnmarshalVT(b)
		vt.unknownFields = nil
		t.Logf("vtprotobuf: %v, %v", protojson.Format(vt), prototext.Format(vt))
		require.Equal(t, us, vt)
		//require.Equal(t, u, v)
	})
}
