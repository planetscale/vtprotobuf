package conformance

import (
	"fmt"
	"testing"
	"time"

	"github.com/planetscale/vtprotobuf/testproto/proto3opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestEquaVTNonNilButOtherwiseEmptyMessages(t *testing.T) {
	plusPrint := func(x *TestAllTypesProto2) string { return fmt.Sprintf("%+v", x) }
	poundPrint := func(x *TestAllTypesProto2) string { return fmt.Sprintf("%#v", x) }
	var same bool

	a := &TestAllTypesProto2{MapStringNestedMessage: map[string]*TestAllTypesProto2_NestedMessage{"": {}}}
	b := &TestAllTypesProto2{MapStringNestedMessage: map[string]*TestAllTypesProto2_NestedMessage{"": {}}}
	require.Equal(t, plusPrint(a), plusPrint(b))
	require.NotEqual(t, poundPrint(a), poundPrint(b))
	same = a.EqualVT(b)
	require.True(t, same)

	c := proto.Clone(a).(*TestAllTypesProto2)
	require.Equal(t, plusPrint(a), plusPrint(c))
	require.NotEqual(t, poundPrint(a), poundPrint(c))
	same = a.EqualVT(c)
	require.True(t, same)
	same = b.EqualVT(c)
	require.True(t, same)

	d := &TestAllTypesProto2{
		MapStringNestedMessage: map[string]*TestAllTypesProto2_NestedMessage{"": (*TestAllTypesProto2_NestedMessage)(nil)}, //
	}
	require.Equal(t, plusPrint(a), plusPrint(d))
	require.NotEqual(t, poundPrint(a), poundPrint(d))
	same = a.EqualVT(d)
	require.True(t, same)
	same = b.EqualVT(d)
	require.True(t, same)
	same = c.EqualVT(d)
	require.True(t, same)
}

func TestEqualVT2(t *testing.T) {
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
			original := proto.Clone(msg).(*TestAllTypesProto2)

			msgData, err := protojson.Marshal(msg)
			require.NoError(t, err)
			originalData, err := protojson.Marshal(original)
			require.NoError(t, err)

			eq := interface{}(msg).(interface {
				EqualVT(*TestAllTypesProto2) bool
			})
			if !eq.EqualVT(original) {
				assert.JSONEq(t, string(originalData), string(msgData))
				err := fmt.Errorf("msg %#v is not EqualVT() to itself %#v", msg, original)
				require.NoError(t, err)
			}

			MutateFields(msg)

			msgData, err = protojson.Marshal(msg)
			require.NoError(t, err)
			originalData, err = protojson.Marshal(original)
			require.NoError(t, err)

			if original.EqualVT(msg) || msg.EqualVT(original) {
				assert.JSONEq(t, string(originalData), string(msgData))
				err = fmt.Errorf("these %T should not be equal:\nmsg = %+v\noriginal = %+v", msg, msg, original)
				require.NoError(t, err)
			}
		})
	}
}

func TestEqualVT3(t *testing.T) {
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

		OptionalDuration:  durationpb.New(time.Hour),
		OptionalTimestamp: timestamppb.Now(),
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

	original := proto.Clone(msg).(*TestAllTypesProto3)

	msgData, err := protojson.Marshal(msg)
	require.NoError(t, err)
	originalData, err := protojson.Marshal(original)
	require.NoError(t, err)

	eq := interface{}(msg).(interface {
		EqualVT(*TestAllTypesProto3) bool
	})
	if !eq.EqualVT(original) {
		assert.JSONEq(t, string(originalData), string(msgData))
		err := fmt.Errorf("msg %#v is not EqualVT() to itself %#v", msg, original)
		require.NoError(t, err)
	}

	MutateFields(msg)

	msgData, err = protojson.Marshal(msg)
	require.NoError(t, err)
	originalData, err = protojson.Marshal(original)
	require.NoError(t, err)

	if original.EqualVT(msg) || msg.EqualVT(original) {
		assert.JSONEq(t, string(originalData), string(msgData))
		err = fmt.Errorf("these %T should not be equal:\nmsg = %+v\noriginal = %+v", msg, msg, original)
		require.NoError(t, err)
	}
}

func TestEqualVT_Map_AbsenceVsZeroValue(t *testing.T) {
	a := &TestAllTypesProto3{
		MapInt32Int32: map[int32]int32{
			1: 0,
			2: 37,
		},
	}
	b := &TestAllTypesProto3{
		MapInt32Int32: map[int32]int32{
			2: 37,
			3: 42,
		},
	}

	aJson, err := protojson.Marshal(a)
	require.NoError(t, err)
	bJson, err := protojson.Marshal(b)
	require.NoError(t, err)

	if a.EqualVT(b) {
		assert.JSONEq(t, string(aJson), string(bJson))
		err := fmt.Errorf("these %T should not be equal:\nmsg = %+v\noriginal = %+v", a, a, b)
		require.NoError(t, err)
	}
}

func TestEqualVT_Oneof_AbsenceVsZeroValue(t *testing.T) {
	a := &TestAllTypesProto3{
		OneofField: &TestAllTypesProto3_OneofUint32{
			OneofUint32: 0,
		},
	}
	b := &TestAllTypesProto3{
		OneofField: &TestAllTypesProto3_OneofString{
			OneofString: "",
		},
	}

	aJson, err := protojson.Marshal(a)
	require.NoError(t, err)
	bJson, err := protojson.Marshal(b)
	require.NoError(t, err)

	if a.EqualVT(b) {
		assert.JSONEq(t, string(aJson), string(bJson))
		err := fmt.Errorf("these %T should not be equal:\nmsg = %+v\noriginal = %+v", a, a, b)
		require.NoError(t, err)
	}
}

func TestEqualVT_Proto2_BytesPresence(t *testing.T) {
	a := &TestAllTypesProto2{
		OptionalBytes: nil,
	}
	b := &TestAllTypesProto2{
		OptionalBytes: []byte{},
	}

	require.False(t, proto.Equal(a, b))

	aJson, err := protojson.Marshal(a)
	require.NoError(t, err)
	bJson, err := protojson.Marshal(b)
	require.NoError(t, err)

	if a.EqualVT(b) {
		assert.JSONEq(t, string(aJson), string(bJson))
		err := fmt.Errorf("these %T should not be equal:\nmsg = %+v\noriginal = %+v", a, a, b)
		require.NoError(t, err)
	}
}

func TestEqualVT_Proto3_BytesPresence(t *testing.T) {
	a := &proto3opt.OptionalFieldInProto3{
		OptionalBytes: nil,
	}
	b := &proto3opt.OptionalFieldInProto3{
		OptionalBytes: []byte{},
	}

	require.False(t, proto.Equal(a, b))

	aJson, err := protojson.Marshal(a)
	require.NoError(t, err)
	bJson, err := protojson.Marshal(b)
	require.NoError(t, err)

	if a.EqualVT(b) {
		assert.JSONEq(t, string(aJson), string(bJson))
		err := fmt.Errorf("these %T should not be equal:\nmsg = %+v\noriginal = %+v", a, a, b)
		require.NoError(t, err)
	}
}

func TestEqualVT_Proto2_BytesNoPresence(t *testing.T) {
	a := &TestAllTypesProto2{
		RepeatedBytes: [][]byte{nil},
		OneofField: &TestAllTypesProto2_OneofBytes{
			OneofBytes: nil,
		},
	}
	b := &TestAllTypesProto2{
		RepeatedBytes: [][]byte{{}},
		OneofField: &TestAllTypesProto2_OneofBytes{
			OneofBytes: []byte{},
		},
	}

	require.True(t, proto.Equal(a, b))

	if !a.EqualVT(b) {
		err := fmt.Errorf("these %T should be equal:\nmsg = %+v\noriginal = %+v", a, a, b)
		require.NoError(t, err)
	}
}

func TestEqualVT_Proto3_BytesNoPresence(t *testing.T) {
	a := &TestAllTypesProto3{
		RepeatedBytes: [][]byte{nil},
		OneofField: &TestAllTypesProto3_OneofBytes{
			OneofBytes: nil,
		},
		OptionalBytes: nil,
	}
	b := &TestAllTypesProto3{
		RepeatedBytes: [][]byte{{}},
		OneofField: &TestAllTypesProto3_OneofBytes{
			OneofBytes: []byte{},
		},
		OptionalBytes: []byte{},
	}

	require.True(t, proto.Equal(a, b))

	if !a.EqualVT(b) {
		err := fmt.Errorf("these %T should not be equal:\nmsg = %+v\noriginal = %+v", a, a, b)
		require.NoError(t, err)
	}
}

func TestEqualVT_NilVsEmpty(t *testing.T) {
	cases := map[string][2]*TestAllTypesProto3{
		"nil and empty should not be equal": {
			&TestAllTypesProto3{},
			(*TestAllTypesProto3)(nil),
		},
		"nil and empty message field should not be equal": {
			&TestAllTypesProto3{
				OptionalNestedMessage: &TestAllTypesProto3_NestedMessage{},
			},
			&TestAllTypesProto3{
				OptionalNestedMessage: nil,
			},
		},
		"nil and empty message should be equal in slice": {
			&TestAllTypesProto3{
				RepeatedNestedMessage: []*TestAllTypesProto3_NestedMessage{{}},
			},
			&TestAllTypesProto3{
				RepeatedNestedMessage: []*TestAllTypesProto3_NestedMessage{nil},
			},
		},
		"nil and empty message should be equal in map value": {
			&TestAllTypesProto3{
				MapStringNestedMessage: map[string]*TestAllTypesProto3_NestedMessage{
					"": {},
				},
			},
			&TestAllTypesProto3{
				MapStringNestedMessage: map[string]*TestAllTypesProto3_NestedMessage{
					"": nil,
				},
			},
		},
		"nil and empty message should be equal in oneof": {
			&TestAllTypesProto3{
				OneofField: &TestAllTypesProto3_OneofNestedMessage{
					OneofNestedMessage: &TestAllTypesProto3_NestedMessage{},
				},
			},
			&TestAllTypesProto3{
				OneofField: &TestAllTypesProto3_OneofNestedMessage{
					OneofNestedMessage: nil,
				},
			},
		},
	}

	for name, c := range cases {
		cc := c // avoid loop closure bug
		t.Run(name, func(t *testing.T) {
			if proto.Equal(cc[0], cc[1]) {
				assert.Truef(t, cc[0].EqualVT(cc[1]), "these %T should be equal:\nfirst = %+v\nsecond = %+v\n", cc[0], cc[0], cc[1])
			} else {
				assert.Falsef(t, cc[0].EqualVT(cc[1]), "these %T should not be equal:\nfirst = %+v\nsecond = %+v\n", cc[0], cc[0], cc[1])
			}
		})
	}
}
