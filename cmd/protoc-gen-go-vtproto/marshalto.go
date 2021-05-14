package main

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type counter int

func (this *counter) Next() string {
	(*this)++
	return this.Current()
}

func (this *counter) Current() string {
	return strconv.Itoa(int(*this))
}

func (p *vtproto) callFixed64(varName ...string) {
	p.P(`i -= 8`)
	p.P(p.Ident("encoding/binary", "LittleEndian"), `.PutUint64(dAtA[i:], uint64(`, strings.Join(varName, ""), `))`)
}

func (p *vtproto) callFixed32(varName ...string) {
	p.P(`i -= 4`)
	p.P(p.Ident("encoding/binary", "LittleEndian"), `.PutUint32(dAtA[i:], uint32(`, strings.Join(varName, ""), `))`)
}

func (p *vtproto) callVarint(varName ...string) {
	p.P(`i = encodeVarint`, p.localName, `(dAtA, i, uint64(`, strings.Join(varName, ""), `))`)
}

func (p *vtproto) encodeKey(fieldNumber protoreflect.FieldNumber, wireType protowire.Type) {
	x := uint32(fieldNumber)<<3 | uint32(wireType)
	i := 0
	keybuf := make([]byte, 0)
	for i = 0; x > 127; i++ {
		keybuf = append(keybuf, 0x80|uint8(x&0x7F))
		x >>= 7
	}
	keybuf = append(keybuf, uint8(x))
	for i = len(keybuf) - 1; i >= 0; i-- {
		p.P(`i--`)
		p.P(`dAtA[i] = `, fmt.Sprintf("%#v", keybuf[i]))
	}
}

func keySize(fieldNumber protoreflect.FieldNumber, wireType protowire.Type) int {
	x := uint32(fieldNumber)<<3 | uint32(wireType)
	size := 0
	for size = 0; x > 127; size++ {
		x >>= 7
	}
	size++
	return size
}

func (p *vtproto) marshalMapField(field *protogen.Field, kvField *protogen.Field, varName string) {
	switch kvField.Desc.Kind() {
	case protoreflect.DoubleKind:
		p.callFixed64(p.Ident("math", "Float64bits"), `(float64(`, varName, `))`)
	case protoreflect.FloatKind:
		p.callFixed32(p.Ident("math", "Float32bits"), `(float32(`, varName, `))`)
	case protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Int32Kind, protoreflect.Uint32Kind, protoreflect.EnumKind:
		p.callVarint(varName)
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		p.callFixed64(varName)
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		p.callFixed32(varName)
	case protoreflect.BoolKind:
		p.P(`i--`)
		p.P(`if `, varName, ` {`)
		p.P(`dAtA[i] = 1`)
		p.P(`} else {`)
		p.P(`dAtA[i] = 0`)
		p.P(`}`)
	case protoreflect.StringKind, protoreflect.BytesKind:
		p.P(`i -= len(`, varName, `)`)
		p.P(`copy(dAtA[i:], `, varName, `)`)
		p.callVarint(`len(`, varName, `)`)
	case protoreflect.Sint32Kind:
		p.callVarint(`(uint32(`, varName, `) << 1) ^ uint32((`, varName, ` >> 31))`)
	case protoreflect.Sint64Kind:
		p.callVarint(`(uint64(`, varName, `) << 1) ^ uint64((`, varName, ` >> 63))`)
	case protoreflect.MessageKind:
		p.marshalBackward(varName, true, kvField.Message)
	}
}

func (p *vtproto) marshalField(proto3 bool, numGen *counter, message *protogen.Message, field *protogen.Field) {
	fieldname := field.GoName
	nullcheck := field.Message != nil
	repeated := field.Desc.Cardinality() == protoreflect.Repeated
	if repeated {
		p.P(`if len(m.`, fieldname, `) > 0 {`)
	} else if nullcheck {
		p.P(`if m.`, fieldname, ` != nil {`)
	}
	packed := field.Desc.IsPacked()
	wireType := wireTypes[field.Desc.Kind()]
	fieldNumber := field.Desc.Number()
	if packed {
		wireType = protowire.BytesType
	}
	switch field.Desc.Kind() {
	case protoreflect.DoubleKind:
		if packed {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`f`, numGen.Next(), ` := `, p.Ident("math", "Float64bits"), `(float64(`, val, `))`)
			p.callFixed64("f" + numGen.Current())
			p.P(`}`)
			p.callVarint(`len(m.`, fieldname, `) * 8`)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`f`, numGen.Next(), ` := `, p.Ident("math", "Float64bits"), `(float64(`, val, `))`)
			p.callFixed64("f" + numGen.Current())
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.callFixed64(p.Ident("math", "Float64bits"), `(float64(m.`, fieldname, `))`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.callFixed64(p.Ident("math", "Float64bits"), `(float64(m.`+fieldname, `))`)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.FloatKind:
		if packed {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`f`, numGen.Next(), ` := `, p.Ident("math", "Float32bits"), `(float32(`, val, `))`)
			p.callFixed32("f" + numGen.Current())
			p.P(`}`)
			p.callVarint(`len(m.`, fieldname, `) * 4`)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`f`, numGen.Next(), ` := `, p.Ident("math", "Float32bits"), `(float32(`, val, `))`)
			p.callFixed32("f" + numGen.Current())
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.callFixed32(p.Ident("math", "Float32bits"), `(float32(m.`+fieldname, `))`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.callFixed32(p.Ident("math", "Float32bits"), `(float32(m.`+fieldname, `))`)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Int32Kind, protoreflect.Uint32Kind, protoreflect.EnumKind:
		if packed {
			jvar := "j" + numGen.Next()
			p.P(`dAtA`, numGen.Next(), ` := make([]byte, len(m.`, fieldname, `)*10)`)
			p.P(`var `, jvar, ` int`)
			switch field.Desc.Kind() {
			case protoreflect.Int64Kind, protoreflect.Int32Kind:
				p.P(`for _, num1 := range m.`, fieldname, ` {`)
				p.P(`num := uint64(num1)`)
			default:
				p.P(`for _, num := range m.`, fieldname, ` {`)
			}
			p.P(`for num >= 1<<7 {`)
			p.P(`dAtA`, numGen.Current(), `[`, jvar, `] = uint8(uint64(num)&0x7f|0x80)`)
			p.P(`num >>= 7`)
			p.P(jvar, `++`)
			p.P(`}`)
			p.P(`dAtA`, numGen.Current(), `[`, jvar, `] = uint8(num)`)
			p.P(jvar, `++`)
			p.P(`}`)
			p.P(`i -= `, jvar)
			p.P(`copy(dAtA[i:], dAtA`, numGen.Current(), `[:`, jvar, `])`)
			p.callVarint(jvar)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.callVarint(val)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.callVarint(`m.`, fieldname)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.callVarint(`m.`, fieldname)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		if packed {
			val := p.reverseListRange(`m.`, fieldname)
			p.callFixed64(val)
			p.P(`}`)
			p.callVarint(`len(m.`, fieldname, `) * 8`)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.callFixed64(val)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.callFixed64("m." + fieldname)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.callFixed64("m." + fieldname)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		if packed {
			val := p.reverseListRange(`m.`, fieldname)
			p.callFixed32(val)
			p.P(`}`)
			p.callVarint(`len(m.`, fieldname, `) * 4`)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.callFixed32(val)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.callFixed32("m." + fieldname)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.callFixed32("m." + fieldname)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.BoolKind:
		if packed {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`i--`)
			p.P(`if `, val, ` {`)
			p.P(`dAtA[i] = 1`)
			p.P(`} else {`)
			p.P(`dAtA[i] = 0`)
			p.P(`}`)
			p.P(`}`)
			p.callVarint(`len(m.`, fieldname, `)`)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`i--`)
			p.P(`if `, val, ` {`)
			p.P(`dAtA[i] = 1`)
			p.P(`} else {`)
			p.P(`dAtA[i] = 0`)
			p.P(`}`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` {`)
			p.P(`i--`)
			p.P(`if m.`, fieldname, ` {`)
			p.P(`dAtA[i] = 1`)
			p.P(`} else {`)
			p.P(`dAtA[i] = 0`)
			p.P(`}`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.P(`i--`)
			p.P(`if m.`, fieldname, ` {`)
			p.P(`dAtA[i] = 1`)
			p.P(`} else {`)
			p.P(`dAtA[i] = 0`)
			p.P(`}`)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.StringKind:
		if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`i -= len(`, val, `)`)
			p.P(`copy(dAtA[i:], `, val, `)`)
			p.callVarint(`len(`, val, `)`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if len(m.`, fieldname, `) > 0 {`)
			p.P(`i -= len(m.`, fieldname, `)`)
			p.P(`copy(dAtA[i:], m.`, fieldname, `)`)
			p.callVarint(`len(m.`, fieldname, `)`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.P(`i -= len(m.`, fieldname, `)`)
			p.P(`copy(dAtA[i:], m.`, fieldname, `)`)
			p.callVarint(`len(m.`, fieldname, `)`)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.GroupKind:
		panic(fmt.Errorf("marshaler does not support group %v", fieldname))
	case protoreflect.MessageKind:
		if field.Desc.IsMap() {
			goTypK, _ := fieldGoType(p.GeneratedFile, field.Message.Fields[0])
			keyKind := field.Message.Fields[0].Desc.Kind()
			valKind := field.Message.Fields[1].Desc.Kind()

			var val string
			if p.stable {
				keysName := `keysFor` + fieldname
				p.P(keysName, ` := make([]`, goTypK, `, 0, len(m.`, fieldname, `))`)
				p.P(`for k := range m.`, fieldname, ` {`)
				p.P(keysName, ` = append(`, keysName, `, `, goTypK, `(k))`)
				p.P(`}`)
				p.P(p.Ident("sort", CamelCase(goTypK)+"s"), `(`, keysName, `)`)
				val = p.reverseListRange(keysName)
			} else {
				p.P(`for k := range m.`, fieldname, ` {`)
				val = "k"
			}
			if p.stable {
				p.P(`v := m.`, fieldname, `[`, goTypK, `(`, val, `)]`)
			} else {
				p.P(`v := m.`, fieldname, `[`, val, `]`)
			}
			p.P(`baseI := i`)

			accessor := `v`
			if valKind == protoreflect.BytesKind {
				if proto3 {
					p.P(`if len(`, accessor, `) > 0 {`)
				} else {
					p.P(`if `, accessor, ` != nil {`)
				}
			}
			p.marshalMapField(field, field.Message.Fields[1], accessor)
			p.encodeKey(2, wireTypes[valKind])
			if valKind == protoreflect.BytesKind {
				p.P(`}`)
			}
			p.marshalMapField(field, field.Message.Fields[0], val)
			p.encodeKey(1, wireTypes[keyKind])
			p.callVarint(`baseI - i`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.marshalBackward(val, true, field.Message)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.marshalBackward(`m.`+fieldname, true, field.Message)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.BytesKind:
		if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`i -= len(`, val, `)`)
			p.P(`copy(dAtA[i:], `, val, `)`)
			p.callVarint(`len(`, val, `)`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if len(m.`, fieldname, `) > 0 {`)
			p.P(`i -= len(m.`, fieldname, `)`)
			p.P(`copy(dAtA[i:], m.`, fieldname, `)`)
			p.callVarint(`len(m.`, fieldname, `)`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.P(`i -= len(m.`, fieldname, `)`)
			p.P(`copy(dAtA[i:], m.`, fieldname, `)`)
			p.callVarint(`len(m.`, fieldname, `)`)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.Sint32Kind:
		if packed {
			datavar := "dAtA" + numGen.Next()
			jvar := "j" + numGen.Next()
			p.P(datavar, ` := make([]byte, len(m.`, fieldname, ")*5)")
			p.P(`var `, jvar, ` int`)
			p.P(`for _, num := range m.`, fieldname, ` {`)
			xvar := "x" + numGen.Next()
			p.P(xvar, ` := (uint32(num) << 1) ^ uint32((num >> 31))`)
			p.P(`for `, xvar, ` >= 1<<7 {`)
			p.P(datavar, `[`, jvar, `] = uint8(uint64(`, xvar, `)&0x7f|0x80)`)
			p.P(jvar, `++`)
			p.P(xvar, ` >>= 7`)
			p.P(`}`)
			p.P(datavar, `[`, jvar, `] = uint8(`, xvar, `)`)
			p.P(jvar, `++`)
			p.P(`}`)
			p.P(`i -= `, jvar)
			p.P(`copy(dAtA[i:], `, datavar, `[:`, jvar, `])`)
			p.callVarint(jvar)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`x`, numGen.Next(), ` := (uint32(`, val, `) << 1) ^ uint32((`, val, ` >> 31))`)
			p.callVarint(`x`, numGen.Current())
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.callVarint(`(uint32(m.`, fieldname, `) << 1) ^ uint32((m.`, fieldname, ` >> 31))`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.callVarint(`(uint32(m.`, fieldname, `) << 1) ^ uint32((m.`, fieldname, ` >> 31))`)
			p.encodeKey(fieldNumber, wireType)
		}
	case protoreflect.Sint64Kind:
		if packed {
			jvar := "j" + numGen.Next()
			xvar := "x" + numGen.Next()
			datavar := "dAtA" + numGen.Next()
			p.P(`var `, jvar, ` int`)
			p.P(datavar, ` := make([]byte, len(m.`, fieldname, `)*10)`)
			p.P(`for _, num := range m.`, fieldname, ` {`)
			p.P(xvar, ` := (uint64(num) << 1) ^ uint64((num >> 63))`)
			p.P(`for `, xvar, ` >= 1<<7 {`)
			p.P(datavar, `[`, jvar, `] = uint8(uint64(`, xvar, `)&0x7f|0x80)`)
			p.P(jvar, `++`)
			p.P(xvar, ` >>= 7`)
			p.P(`}`)
			p.P(datavar, `[`, jvar, `] = uint8(`, xvar, `)`)
			p.P(jvar, `++`)
			p.P(`}`)
			p.P(`i -= `, jvar)
			p.P(`copy(dAtA[i:], `, datavar, `[:`, jvar, `])`)
			p.callVarint(jvar)
			p.encodeKey(fieldNumber, wireType)
		} else if repeated {
			val := p.reverseListRange(`m.`, fieldname)
			p.P(`x`, numGen.Next(), ` := (uint64(`, val, `) << 1) ^ uint64((`, val, ` >> 63))`)
			p.callVarint("x" + numGen.Current())
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.callVarint(`(uint64(m.`, fieldname, `) << 1) ^ uint64((m.`, fieldname, ` >> 63))`)
			p.encodeKey(fieldNumber, wireType)
			p.P(`}`)
		} else {
			p.callVarint(`(uint64(m.`, fieldname, `) << 1) ^ uint64((m.`, fieldname, ` >> 63))`)
			p.encodeKey(fieldNumber, wireType)
		}
	default:
		panic("not implemented")
	}
	if repeated || nullcheck {
		p.P(`}`)
	}
}

func (p *vtproto) generateMessageMarshal(message *protogen.Message) {
	var numGen counter
	ccTypeName := message.GoIdent

	p.P(`func (m *`, ccTypeName, `) MarshalVT() (dAtA []byte, err error) {`)
	p.P(`if m == nil {`)
	p.P(`return nil, nil`)
	p.P(`}`)
	p.P(`size := m.SizeVT()`)
	p.P(`dAtA = make([]byte, size)`)
	p.P(`n, err := m.MarshalToSizedBufferVT(dAtA[:size])`)
	p.P(`if err != nil {`)
	p.P(`return nil, err`)
	p.P(`}`)
	p.P(`return dAtA[:n], nil`)
	p.P(`}`)
	p.P(``)
	p.P(`func (m *`, ccTypeName, `) MarshalToVT(dAtA []byte) (int, error) {`)
	p.P(`size := m.SizeVT()`)
	p.P(`return m.MarshalToSizedBufferVT(dAtA[:size])`)
	p.P(`}`)
	p.P(``)
	p.P(`func (m *`, ccTypeName, `) MarshalToSizedBufferVT(dAtA []byte) (int, error) {`)
	// p.P(`if m == nil {`)
	// p.P(`return 0, nil`)
	// p.P(`}`)
	p.P(`i := len(dAtA)`)
	p.P(`_ = i`)
	p.P(`var l int`)
	p.P(`_ = l`)
	p.P(`if m.unknownFields != nil {`)
	p.P(`i -= len(m.unknownFields)`)
	p.P(`copy(dAtA[i:], m.unknownFields)`)
	p.P(`}`)
	oneofs := make(map[string]struct{})
	for i := len(message.Fields) - 1; i >= 0; i-- {
		field := message.Fields[i]
		oneof := field.Oneof != nil
		if !oneof {
			p.marshalField(true, &numGen, message, field)
		} else {
			fieldname := field.Oneof.GoName
			if _, ok := oneofs[fieldname]; !ok {
				oneofs[fieldname] = struct{}{}
				p.P(`if vtmsg, ok := m.`, fieldname, `.(vtprotoMessage`, p.localName, `); ok {`)
				p.marshalForward("vtmsg", false, false)
				p.P(`}`)
			}
		}
	}
	p.P(`return len(dAtA) - i, nil`)
	p.P(`}`)
	p.P()

	//Generate MarshalToVT methods for oneof fields
	for _, field := range message.Fields {
		if field.Oneof == nil {
			continue
		}
		ccTypeName := field.GoIdent
		p.P(`func (m *`, ccTypeName, `) MarshalToVT(dAtA []byte) (int, error) {`)
		p.P(`size := m.SizeVT()`)
		p.P(`return m.MarshalToSizedBufferVT(dAtA[:size])`)
		p.P(`}`)
		p.P(``)
		p.P(`func (m *`, ccTypeName, `) MarshalToSizedBufferVT(dAtA []byte) (int, error) {`)
		p.P(`i := len(dAtA)`)
		p.marshalField(false, &numGen, message, field)
		p.P(`return len(dAtA) - i, nil`)
		p.P(`}`)
	}
}

func (p *vtproto) generateMarshalHelpers() {
	p.P(`func encodeVarint`, p.localName, `(dAtA []byte, offset int, v uint64) int {`)
	p.P(`offset -= sov`, p.localName, `(v)`)
	p.P(`base := offset`)
	p.P(`for v >= 1<<7 {`)
	p.P(`dAtA[offset] = uint8(v&0x7f|0x80)`)
	p.P(`v >>= 7`)
	p.P(`offset++`)
	p.P(`}`)
	p.P(`dAtA[offset] = uint8(v)`)
	p.P(`return base`)
	p.P(`}`)
}

func (p *vtproto) reverseListRange(expression ...string) string {
	exp := strings.Join(expression, "")
	p.P(`for iNdEx := len(`, exp, `) - 1; iNdEx >= 0; iNdEx-- {`)
	return exp + `[iNdEx]`
}

func (p *vtproto) marshalBackward(varName string, varInt bool, message *protogen.Message) {
	foreign := strings.HasPrefix(string(message.Desc.FullName()), "google.protobuf.")

	p.P(`{`)
	if foreign {
		p.P(`encoded, err := `, p.Ident(ProtoPkg, "Marshal"), `(`, varName, `)`)
	} else {
		p.P(`size, err := `, varName, `.MarshalToSizedBufferVT(dAtA[:i])`)
	}

	p.P(`if err != nil {`)
	p.P(`return 0, err`)
	p.P(`}`)

	if foreign {
		p.P(`i -= len(encoded)`)
		p.P(`copy(dAtA[i:], encoded)`)
		if varInt {
			p.callVarint(`len(encoded)`)
		}
	} else {
		p.P(`i -= size`)
		if varInt {
			p.callVarint(`size`)
		}
	}
	p.P(`}`)
}

func (p *vtproto) marshalForward(varName string, varInt, protoSizer bool) {
	p.P(`{`)
	p.P(`size := `, varName, `.SizeVT()`)
	p.P(`i -= size`)
	p.P(`if _, err := `, varName, `.MarshalToVT(dAtA[i:]); err != nil {`)
	p.P(`return 0, err`)
	p.P(`}`)
	if varInt {
		p.callVarint(`size`)
	}
	p.P(`}`)
}
