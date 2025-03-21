// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: unique/unique.proto

package unique

import (
	_ "github.com/planetscale/vtprotobuf/vtproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UniqueFieldExtension struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Foo string           `protobuf:"bytes,1,opt,name=foo,proto3" json:"foo,omitempty"`
	Bar map[string]int64 `protobuf:"bytes,2,rep,name=bar,proto3" json:"bar,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Baz map[int64]string `protobuf:"bytes,3,rep,name=baz,proto3" json:"baz,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *UniqueFieldExtension) Reset() {
	*x = UniqueFieldExtension{}
	if protoimpl.UnsafeEnabled {
		mi := &file_unique_unique_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UniqueFieldExtension) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UniqueFieldExtension) ProtoMessage() {}

func (x *UniqueFieldExtension) ProtoReflect() protoreflect.Message {
	mi := &file_unique_unique_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UniqueFieldExtension.ProtoReflect.Descriptor instead.
func (*UniqueFieldExtension) Descriptor() ([]byte, []int) {
	return file_unique_unique_proto_rawDescGZIP(), []int{0}
}

func (x *UniqueFieldExtension) GetFoo() string {
	if x != nil {
		return x.Foo
	}
	return ""
}

func (x *UniqueFieldExtension) GetBar() map[string]int64 {
	if x != nil {
		return x.Bar
	}
	return nil
}

func (x *UniqueFieldExtension) GetBaz() map[int64]string {
	if x != nil {
		return x.Baz
	}
	return nil
}

var File_unique_unique_proto protoreflect.FileDescriptor

var file_unique_unique_proto_rawDesc = []byte{
	0x0a, 0x13, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x2f, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x2f, 0x76, 0x74,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x65, 0x78, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x02, 0x0a, 0x14, 0x55,
	0x6e, 0x69, 0x71, 0x75, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x03, 0x66, 0x6f, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x06, 0xb2, 0xa9, 0x1f, 0x02, 0x08, 0x01, 0x52, 0x03, 0x66, 0x6f, 0x6f, 0x12, 0x38, 0x0a,
	0x03, 0x62, 0x61, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x55, 0x6e, 0x69,
	0x71, 0x75, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x42, 0x61, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x06, 0xb2, 0xa9, 0x1f, 0x02,
	0x08, 0x01, 0x52, 0x03, 0x62, 0x61, 0x72, 0x12, 0x38, 0x0a, 0x03, 0x62, 0x61, 0x7a, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x42, 0x61, 0x7a, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x42, 0x06, 0xb2, 0xa9, 0x1f, 0x02, 0x08, 0x01, 0x52, 0x03, 0x62, 0x61,
	0x7a, 0x1a, 0x36, 0x0a, 0x08, 0x42, 0x61, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x36, 0x0a, 0x08, 0x42, 0x61, 0x7a,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x42, 0x12, 0x5a, 0x10, 0x74, 0x65, 0x73, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x75,
	0x6e, 0x69, 0x71, 0x75, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_unique_unique_proto_rawDescOnce sync.Once
	file_unique_unique_proto_rawDescData = file_unique_unique_proto_rawDesc
)

func file_unique_unique_proto_rawDescGZIP() []byte {
	file_unique_unique_proto_rawDescOnce.Do(func() {
		file_unique_unique_proto_rawDescData = protoimpl.X.CompressGZIP(file_unique_unique_proto_rawDescData)
	})
	return file_unique_unique_proto_rawDescData
}

var file_unique_unique_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_unique_unique_proto_goTypes = []interface{}{
	(*UniqueFieldExtension)(nil), // 0: UniqueFieldExtension
	nil,                          // 1: UniqueFieldExtension.BarEntry
	nil,                          // 2: UniqueFieldExtension.BazEntry
}
var file_unique_unique_proto_depIdxs = []int32{
	1, // 0: UniqueFieldExtension.bar:type_name -> UniqueFieldExtension.BarEntry
	2, // 1: UniqueFieldExtension.baz:type_name -> UniqueFieldExtension.BazEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_unique_unique_proto_init() }
func file_unique_unique_proto_init() {
	if File_unique_unique_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_unique_unique_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UniqueFieldExtension); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_unique_unique_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_unique_unique_proto_goTypes,
		DependencyIndexes: file_unique_unique_proto_depIdxs,
		MessageInfos:      file_unique_unique_proto_msgTypes,
	}.Build()
	File_unique_unique_proto = out.File
	file_unique_unique_proto_rawDesc = nil
	file_unique_unique_proto_goTypes = nil
	file_unique_unique_proto_depIdxs = nil
}
