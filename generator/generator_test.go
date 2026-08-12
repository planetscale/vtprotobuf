package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/gofeaturespb"
)

type lazyFieldDescriptor struct {
	protoreflect.FieldDescriptor
}

func (lazyFieldDescriptor) IsLazy() bool { return true }

func TestValidateRejectsWrapWithOpaqueAPI(t *testing.T) {
	msg := testMessageDescriptor(t)
	gen := &Generator{cfg: &Config{Wrap: true}}

	err := gen.validate(&protogen.File{Messages: []*protogen.Message{{
		Desc:     msg,
		APILevel: gofeaturespb.GoFeatures_API_OPAQUE,
	}}})

	require.ErrorContains(t, err, "wrap option is incompatible with the opaque API")
}

func TestValidateRejectsLazyFields(t *testing.T) {
	msg := testMessageDescriptor(t)
	field := msg.Fields().ByName("lazy_field")
	gen := &Generator{cfg: &Config{}}

	err := gen.validate(&protogen.File{Messages: []*protogen.Message{{
		Desc: msg,
		Fields: []*protogen.Field{{
			Desc: lazyFieldDescriptor{FieldDescriptor: field},
		}},
	}}})

	require.ErrorContains(t, err, "lazy fields are not supported")
}

func testMessageDescriptor(t *testing.T) protoreflect.MessageDescriptor {
	t.Helper()

	fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Package: proto.String("testpb"),
		Syntax:  proto.String("proto2"),
		MessageType: []*descriptorpb.DescriptorProto{
			{Name: proto.String("Inner")},
			{
				Name: proto.String("TestMessage"),
				Field: []*descriptorpb.FieldDescriptorProto{{
					Name:     proto.String("lazy_field"),
					Number:   proto.Int32(1),
					Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
					TypeName: proto.String(".testpb.Inner"),
					Options: &descriptorpb.FieldOptions{
						Lazy: proto.Bool(true),
					},
				}},
			},
		},
	}, nil)
	require.NoError(t, err)

	return fd.Messages().ByName("TestMessage")
}
