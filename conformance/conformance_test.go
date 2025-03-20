// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package conformance_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	pb "github.com/planetscale/vtprotobuf/conformance/internal/conformance"
)

func init() {
	// When the environment variable RUN_AS_CONFORMANCE_PLUGIN is set,
	// we skip running the tests and instead act as a conformance plugin.
	// This allows the binary to pass itself to conformance.
	if os.Getenv("RUN_AS_CONFORMANCE_PLUGIN") == "1" {
		main()
		os.Exit(0)
	}
}

var (
	execute   = flag.Bool("execute", true, "execute the conformance test")
	protoRoot = flag.String("protoroot", os.Getenv("PROTOBUF_ROOT"), "The root of the protobuf source tree.")
)

func Test(t *testing.T) {
	if !*execute || testing.Short() {
		t.SkipNow()
	}
	binPath := filepath.Join(*protoRoot, "conformance_test_runner")
	cmd := exec.Command(binPath,
		// "--failure_list", "failing_tests.txt",
		"--text_format_failure_list", "failing_tests_text_format.txt",
		"--enforce_recommended",
		os.Args[0])
	cmd.Env = append(os.Environ(), "RUN_AS_CONFORMANCE_PLUGIN=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("execution error: %v\n\n%s", err, out)
	}
}

var marshalDifflog io.WriteCloser

func conformanceUnmarshal(b []byte, msg proto.Message) error {
	expected := proto.Clone(msg)
	if err := proto.Unmarshal(b, expected); err != nil {
		return err
	}
	type unmarshalvt interface {
		UnmarshalVT(b []byte) error
	}
	u := msg.(unmarshalvt)
	if err := u.UnmarshalVT(b); err != nil {
		return err
	}

	if !proto.Equal(expected, msg) {
		fmt.Fprintf(marshalDifflog, "UNMARSHAL\n")
		fmt.Fprintf(marshalDifflog, "expected:\n%s\n\n", prototext.Format(expected))
		fmt.Fprintf(marshalDifflog, "got:\n%s\n\n", prototext.Format(msg))
		fmt.Fprintf(marshalDifflog, "raw: %#v\n\n", b)
		fmt.Fprintf(marshalDifflog, "==============\n\n")
	}
	return nil
}

func conformanceMarshal(msg proto.Message) ([]byte, error) {
	var expected, got []byte
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(marshalDifflog, "MARSHAL\n")
			fmt.Fprintf(marshalDifflog, "message:\n%s\n\n", prototext.Format(msg))
			fmt.Fprintf(marshalDifflog, "expected:\n%s\n\n", hex.Dump(expected))
			fmt.Fprintf(marshalDifflog, "CRASH:\n%s\n\n", r)
			fmt.Fprintf(marshalDifflog, "golang:\n%#v\n\n", msg)
			fmt.Fprintf(marshalDifflog, "==============\n\n")
		} else if err != nil {
			// do nothing
		} else if got != nil && !bytes.Equal(expected, got) {
			fmt.Fprintf(marshalDifflog, "MARSHAL\n")
			fmt.Fprintf(marshalDifflog, "message:\n%s\n\n", prototext.Format(msg))
			fmt.Fprintf(marshalDifflog, "expected:\n%s\n\n", hex.Dump(expected))
			fmt.Fprintf(marshalDifflog, "got:\n%s\n\n", hex.Dump(got))
			fmt.Fprintf(marshalDifflog, "golang:\n%#v\n\n", msg)
			fmt.Fprintf(marshalDifflog, "==============\n\n")
		}
	}()

	if expected, err = proto.Marshal(msg); err != nil {
		return nil, err
	}

	type marshalvt interface {
		MarshalVT() ([]byte, error)
	}
	m := msg.(marshalvt)
	if got, err = m.MarshalVT(); err != nil {
		return nil, err
	}
	return got, nil
}

func conformanceEqual(msg proto.Message) (err error) {
	isNaN := func(w interface{}) bool {
		f32, ok32 := w.(float32)
		f64, ok64 := w.(float64)
		return (ok32 && math.IsNaN(float64(f32))) || (ok64 && math.IsNaN(f64))
	}

	switch msg := msg.(type) {
	case *pb.TestAllTypesProto2:
		cloned := proto.Clone(msg).(*pb.TestAllTypesProto2)

		eq := interface{}(msg).(interface {
			EqualVT(*pb.TestAllTypesProto2) bool
		})
		if !eq.EqualVT(cloned) {
			return fmt.Errorf("msg %#v is not EqualVT() to itself %#v", msg, cloned)
		}

		pb.MutateFields(cloned)
		if cloned.EqualVT(msg) || msg.EqualVT(cloned) {
			return fmt.Errorf("these %T should not be equal:\nmsg = %+v\ncloned = %+v", msg, msg, cloned)
		}

	case *pb.TestAllTypesProto3:
		cloned := proto.Clone(msg).(*pb.TestAllTypesProto3)

		eq := interface{}(msg).(interface {
			EqualVT(*pb.TestAllTypesProto3) bool
		})
		same := eq.EqualVT(cloned)
		if pb.VisitWithPredicate(msg, isNaN) {
			if same {
				return fmt.Errorf("msg %T %+v contains NaN thus should not EqualVT() to itself %+v", msg, msg, cloned)
			}
		} else {
			if !same {
				return fmt.Errorf("msg %#v is not EqualVT() to itself %#v", msg, cloned)
			}
		}

		pb.MutateFields(cloned)
		if cloned.EqualVT(msg) || msg.EqualVT(cloned) {
			return fmt.Errorf("these %T should not be equal:\nmsg = %+v\ncloned = %+v", msg, msg, cloned)
		}

	default:
		return fmt.Errorf("unhandled %T", msg)
	}
	return nil
}

func main() {
	var err error
	if marshalDifflog, err = os.Create("marshal.log"); err != nil {
		log.Fatalf("failed to init: %v", err)
	}
	defer marshalDifflog.Close()

	var sizeBuf [4]byte
	inbuf := make([]byte, 0, 4096)
	for {
		_, err := io.ReadFull(os.Stdin, sizeBuf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("conformance: read request: %v", err)
		}
		size := binary.LittleEndian.Uint32(sizeBuf[:])
		if int(size) > cap(inbuf) {
			inbuf = make([]byte, size)
		}
		inbuf = inbuf[:size]
		if _, err := io.ReadFull(os.Stdin, inbuf); err != nil {
			log.Fatalf("conformance: read request: %v", err)
		}

		req := &pb.ConformanceRequest{}
		if err := conformanceUnmarshal(inbuf, req); err != nil {
			log.Fatalf("conformance: parse request: %v", err)
		}
		res := handle(req)

		out, err := conformanceMarshal(res)
		if err != nil {
			log.Fatalf("conformance: marshal response: %v", err)
		}
		binary.LittleEndian.PutUint32(sizeBuf[:], uint32(len(out)))
		if _, err := os.Stdout.Write(sizeBuf[:]); err != nil {
			log.Fatalf("conformance: write response: %v", err)
		}
		if _, err := os.Stdout.Write(out); err != nil {
			log.Fatalf("conformance: write response: %v", err)
		}
	}
}

func handle(req *pb.ConformanceRequest) (res *pb.ConformanceResponse) {
	var msg proto.Message = &pb.TestAllTypesProto2{}
	if req.GetMessageType() == "protobuf_test_messages.proto3.TestAllTypesProto3" {
		msg = &pb.TestAllTypesProto3{}
	}

	// Unmarshal the test message.
	var err error
	switch p := req.Payload.(type) {
	case *pb.ConformanceRequest_ProtobufPayload:
		err = conformanceUnmarshal(p.ProtobufPayload, msg)
	case *pb.ConformanceRequest_JsonPayload:
		err = protojson.UnmarshalOptions{
			DiscardUnknown: req.TestCategory == pb.TestCategory_JSON_IGNORE_UNKNOWN_PARSING_TEST,
		}.Unmarshal([]byte(p.JsonPayload), msg)
	case *pb.ConformanceRequest_TextPayload:
		err = prototext.Unmarshal([]byte(p.TextPayload), msg)
	default:
		return &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_RuntimeError{
				RuntimeError: "unknown request payload type",
			},
		}
	}
	if err != nil {
		return &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_ParseError{
				ParseError: err.Error(),
			},
		}
	}

	if err = conformanceEqual(msg); err != nil {
		return &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_RuntimeError{
				RuntimeError: err.Error(),
			},
		}
	}

	// Marshal the test message.
	var b []byte
	switch req.RequestedOutputFormat {
	case pb.WireFormat_PROTOBUF:
		b, err = conformanceMarshal(msg)
		res = &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_ProtobufPayload{
				ProtobufPayload: b,
			},
		}
	case pb.WireFormat_JSON:
		b, err = protojson.Marshal(msg)
		res = &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_JsonPayload{
				JsonPayload: string(b),
			},
		}
	case pb.WireFormat_TEXT_FORMAT:
		b, err = prototext.MarshalOptions{
			EmitUnknown: req.PrintUnknownFields,
		}.Marshal(msg)
		res = &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_TextPayload{
				TextPayload: string(b),
			},
		}
	default:
		return &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_RuntimeError{
				RuntimeError: "unknown output format",
			},
		}
	}
	if err != nil {
		return &pb.ConformanceResponse{
			Result: &pb.ConformanceResponse_SerializeError{
				SerializeError: err.Error(),
			},
		}
	}
	return res
}
