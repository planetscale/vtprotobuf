package conformance

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCloneVTNil(t *testing.T) {
	var a *TestAllTypesProto2
	b := a.CloneVT()
	require.Nil(t, b)
}

func TestCloneVT2(t *testing.T) {
	stringPtr := func(x string) *string { return &x }
	float32Ptr := func(x float32) *float32 { return &x }
	float64Ptr := func(x float64) *float64 { return &x }

	msgs := []*TestAllTypesProto2{
		{OptionalString: stringPtr("bla")},
		{OptionalDouble: float64Ptr(1.7976931348623157e+308)},
		{OptionalFloat: float32Ptr(-0.0)},
		{OptionalBytes: []byte{}},
		{MapStringBytes: map[string][]byte{"": {}}},
		{OneofField: &TestAllTypesProto2_OneofBool{OneofBool: false}},
		{MapStringNestedMessage: map[string]*TestAllTypesProto2_NestedMessage{"": {}}},
		{MapStringNestedMessage: map[string]*TestAllTypesProto2_NestedMessage{"eh": {}}},
	}

	for _, msg := range msgs {
		t.Run(fmt.Sprintf("%+v", msg), func(t *testing.T) {
			orig := proto.Clone(msg).(*TestAllTypesProto2)
			clone := orig.CloneVT()
			require.Truef(t, clone.EqualVT(msg), "cloning %T returned modified message:\nmsg = %+v\nclone = %+v\n", msg, msg, clone)
			require.Truef(t, orig.EqualVT(msg), "cloning %T modified original message:\nmsg = %+v\nafter clone = %+v\n", msg, msg, clone)

			MutateFields(clone)
			require.False(t, clone.EqualVT(msg), "cloned message unchanged after mutation")
			require.True(t, orig.EqualVT(msg), "mutating cloned %T mutated original:\nmsg = %+v\nafter clone = %+v\n", msg, msg, orig)
		})
	}
}

func TestCloneVT3(t *testing.T) {
	msg := &TestAllTypesProto3{
		OneofField: &TestAllTypesProto3_OneofNullValue{OneofNullValue: structpb.NullValue_NULL_VALUE},

		OptionalBoolWrapper:   wrapperspb.Bool(true),
		OptionalInt32Wrapper:  wrapperspb.Int32(1),
		OptionalInt64Wrapper:  wrapperspb.Int64(1),
		OptionalUint32Wrapper: wrapperspb.UInt32(1),
		OptionalUint64Wrapper: wrapperspb.UInt64(1),
		OptionalFloatWrapper:  wrapperspb.Float(1),
		OptionalDoubleWrapper: wrapperspb.Double(1),
		OptionalStringWrapper: wrapperspb.String("blip"),
		OptionalBytesWrapper:  wrapperspb.Bytes([]byte("blop")),

		RepeatedBoolWrapper:   []*wrapperspb.BoolValue{wrapperspb.Bool(true)},
		RepeatedInt32Wrapper:  []*wrapperspb.Int32Value{wrapperspb.Int32(1)},
		RepeatedInt64Wrapper:  []*wrapperspb.Int64Value{wrapperspb.Int64(1)},
		RepeatedUint32Wrapper: []*wrapperspb.UInt32Value{wrapperspb.UInt32(1)},
		RepeatedUint64Wrapper: []*wrapperspb.UInt64Value{wrapperspb.UInt64(1)},
		RepeatedFloatWrapper:  []*wrapperspb.FloatValue{wrapperspb.Float(1)},
		RepeatedDoubleWrapper: []*wrapperspb.DoubleValue{wrapperspb.Double(1)},
		RepeatedStringWrapper: []*wrapperspb.StringValue{wrapperspb.String("blip")},
		RepeatedBytesWrapper:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("blop"))},

		// OptionalDuration:      *durationpb.Duration
		// OptionalTimestamp:     *timestamppb.Timestamp
		// OptionalFieldMask:     *fieldmaskpb.FieldMask
		// OptionalStruct:        *structpb.Struct
		// OptionalAny:           *anypb.Any
		OptionalValue: structpb.NewNumberValue(42),
		// OptionalNullValue:     structpb.NullValue

		// repeated google.protobuf.Duration repeated_duration
		// repeated google.protobuf.Timestamp repeated_timestamp
		// repeated google.protobuf.FieldMask repeated_fieldmask
		// repeated google.protobuf.Struct repeated_struct
		// repeated google.protobuf.Any repeated_any
		// repeated google.protobuf.Value repeated_value
		RepeatedValue: []*structpb.Value{structpb.NewNumberValue(42)},
		// repeated google.protobuf.ListValue repeated_list_value
	}

	orig := proto.Clone(msg).(*TestAllTypesProto3)
	clone := orig.CloneVT()
	require.Truef(t, clone.EqualVT(msg), "cloning %T returned modified message:\nmsg = %+v\nclone = %+v\n", msg, msg, clone)
	require.Truef(t, orig.EqualVT(msg), "cloning %T modified original message:\nmsg = %+v\nafter clone = %+v\n", msg, msg, clone)

	MutateFields(clone)
	require.False(t, clone.EqualVT(msg), "cloned message unchanged after mutation")
	require.True(t, orig.EqualVT(msg), "mutating cloned %T mutated original:\nmsg = %+v\nafter clone = %+v\n", msg, msg, orig)
}
