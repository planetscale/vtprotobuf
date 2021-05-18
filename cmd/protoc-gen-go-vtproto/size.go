package main

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/encoding/protowire"

	"google.golang.org/protobuf/compiler/protogen"
)

func (p *vtprotofile) sizeForField(proto3 bool, field *protogen.Field, sizeName string) {
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
	key := keySize(fieldNumber, wireType)
	switch field.Desc.Kind() {
	case protoreflect.DoubleKind, protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
		if packed {
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(len(m.`, fieldname, `)*8))`, `+len(m.`, fieldname, `)*8`)
		} else if repeated {
			p.P(`n+=`, strconv.Itoa(key+8), `*len(m.`, fieldname, `)`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.P(`n+=`, strconv.Itoa(key+8))
			p.P(`}`)
		} else {
			p.P(`n+=`, strconv.Itoa(key+8))
		}
	case protoreflect.FloatKind, protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
		if packed {
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(len(m.`, fieldname, `)*4))`, `+len(m.`, fieldname, `)*4`)
		} else if repeated {
			p.P(`n+=`, strconv.Itoa(key+4), `*len(m.`, fieldname, `)`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.P(`n+=`, strconv.Itoa(key+4))
			p.P(`}`)
		} else {
			p.P(`n+=`, strconv.Itoa(key+4))
		}
	case protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Uint32Kind, protoreflect.EnumKind, protoreflect.Int32Kind:
		if packed {
			p.P(`l = 0`)
			p.P(`for _, e := range m.`, fieldname, ` {`)
			p.P(`l+=sov(uint64(e))`)
			p.P(`}`)
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(l))+l`)
		} else if repeated {
			p.P(`for _, e := range m.`, fieldname, ` {`)
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(e))`)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(m.`, fieldname, `))`)
			p.P(`}`)
		} else {
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(m.`, fieldname, `))`)
		}
	case protoreflect.BoolKind:
		if packed {
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(len(m.`, fieldname, `)))`, `+len(m.`, fieldname, `)*1`)
		} else if repeated {
			p.P(`n+=`, strconv.Itoa(key+1), `*len(m.`, fieldname, `)`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` {`)
			p.P(`n+=`, strconv.Itoa(key+1))
			p.P(`}`)
		} else {
			p.P(`n+=`, strconv.Itoa(key+1))
		}
	case protoreflect.StringKind:
		if repeated {
			p.P(`for _, s := range m.`, fieldname, ` { `)
			p.P(`l = len(s)`)
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
			p.P(`}`)
		} else if proto3 {
			p.P(`l=len(m.`, fieldname, `)`)
			p.P(`if l > 0 {`)
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
			p.P(`}`)
		} else {
			p.P(`l=len(m.`, fieldname, `)`)
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
		}
	case protoreflect.GroupKind:
		panic(fmt.Errorf("size does not support group %v", fieldname))
	case protoreflect.MessageKind:
		foreign := strings.HasPrefix(string(field.Message.Desc.FullName()), "google.protobuf.")

		if field.Desc.IsMap() {
			fieldKeySize := keySize(field.Desc.Number(), wireTypes[field.Desc.Kind()])
			keyKeySize := keySize(1, wireTypes[field.Message.Fields[0].Desc.Kind()])
			valueKeySize := keySize(2, wireTypes[field.Message.Fields[1].Desc.Kind()])
			p.P(`for k, v := range m.`, fieldname, ` { `)
			p.P(`_ = k`)
			p.P(`_ = v`)
			sum := []string{strconv.Itoa(keyKeySize)}

			switch field.Message.Fields[0].Desc.Kind() {
			case protoreflect.DoubleKind, protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
				sum = append(sum, `8`)
			case protoreflect.FloatKind, protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
				sum = append(sum, `4`)
			case protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Uint32Kind, protoreflect.EnumKind, protoreflect.Int32Kind:
				sum = append(sum, `sov(uint64(k))`)
			case protoreflect.BoolKind:
				sum = append(sum, `1`)
			case protoreflect.StringKind, protoreflect.BytesKind:
				sum = append(sum, `len(k)+sov(uint64(len(k)))`)
			case protoreflect.Sint32Kind, protoreflect.Sint64Kind:
				sum = append(sum, `soz(uint64(k))`)
			}

			switch field.Message.Fields[1].Desc.Kind() {
			case protoreflect.DoubleKind, protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind:
				sum = append(sum, strconv.Itoa(valueKeySize))
				sum = append(sum, strconv.Itoa(8))
			case protoreflect.FloatKind, protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind:
				sum = append(sum, strconv.Itoa(valueKeySize))
				sum = append(sum, strconv.Itoa(4))
			case protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Uint32Kind, protoreflect.EnumKind, protoreflect.Int32Kind:
				sum = append(sum, strconv.Itoa(valueKeySize))
				sum = append(sum, `sov(uint64(v))`)
			case protoreflect.BoolKind:
				sum = append(sum, strconv.Itoa(valueKeySize))
				sum = append(sum, `1`)
			case protoreflect.StringKind:
				sum = append(sum, strconv.Itoa(valueKeySize))
				sum = append(sum, `len(v)+sov(uint64(len(v)))`)
			case protoreflect.BytesKind:
				p.P(`l = `, strconv.Itoa(valueKeySize), ` + len(v)+sov(uint64(len(v)))`)
				sum = append(sum, `l`)
			case protoreflect.Sint32Kind, protoreflect.Sint64Kind:
				sum = append(sum, strconv.Itoa(valueKeySize))
				sum = append(sum, `soz(uint64(v))`)
			case protoreflect.MessageKind:
				p.P(`l = 0`)
				p.P(`if v != nil {`)
				p.P(`l = v.`, sizeName, `()`)
				p.P(`}`)
				p.P(`l += `, strconv.Itoa(valueKeySize), ` + sov(uint64(l))`)
				sum = append(sum, `l`)
			}
			p.P(`mapEntrySize := `, strings.Join(sum, "+"))
			p.P(`n+=mapEntrySize+`, fieldKeySize, `+sov(uint64(mapEntrySize))`)
			p.P(`}`)
		} else if field.Desc.IsList() {
			p.P(`for _, e := range m.`, fieldname, ` { `)
			if foreign {
				p.P(`l=`, p.Ident(ProtoPkg, "Size"), `(e)`)
			} else {
				p.P(`l=e.`, sizeName, `()`)
			}
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
			p.P(`}`)
		} else {
			if foreign {
				p.P(`l=`, p.Ident(ProtoPkg, "Size"), `(m.`, fieldname, `)`)
			} else {
				p.P(`l=m.`, fieldname, `.`, sizeName, `()`)
			}
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
		}
	case protoreflect.BytesKind:
		if repeated {
			p.P(`for _, b := range m.`, fieldname, ` { `)
			p.P(`l = len(b)`)
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
			p.P(`}`)
		} else if proto3 {
			p.P(`l=len(m.`, fieldname, `)`)
			p.P(`if l > 0 {`)
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
			p.P(`}`)
		} else {
			p.P(`l=len(m.`, fieldname, `)`)
			p.P(`n+=`, strconv.Itoa(key), `+l+sov(uint64(l))`)
		}
	case protoreflect.Sint32Kind, protoreflect.Sint64Kind:
		if packed {
			p.P(`l = 0`)
			p.P(`for _, e := range m.`, fieldname, ` {`)
			p.P(`l+=soz(uint64(e))`)
			p.P(`}`)
			p.P(`n+=`, strconv.Itoa(key), `+sov(uint64(l))+l`)
		} else if repeated {
			p.P(`for _, e := range m.`, fieldname, ` {`)
			p.P(`n+=`, strconv.Itoa(key), `+soz(uint64(e))`)
			p.P(`}`)
		} else if proto3 {
			p.P(`if m.`, fieldname, ` != 0 {`)
			p.P(`n+=`, strconv.Itoa(key), `+soz(uint64(m.`, fieldname, `))`)
			p.P(`}`)
		} else {
			p.P(`n+=`, strconv.Itoa(key), `+soz(uint64(m.`, fieldname, `))`)
		}
	default:
		panic("not implemented")
	}
	if repeated || nullcheck {
		p.P(`}`)
	}
}

func (p *vtprotofile) SizeMessage(message *protogen.Message) {
	sizeName := "SizeVT"
	ccTypeName := message.GoIdent

	p.P(`func (m *`, ccTypeName, `) `, sizeName, `() (n int) {`)
	p.P(`if m == nil {`)
	p.P(`return 0`)
	p.P(`}`)
	p.P(`var l int`)
	p.P(`_ = l`)
	oneofs := make(map[string]struct{})
	for _, field := range message.Fields {
		oneof := field.Oneof != nil
		if !oneof {
			p.sizeForField(true, field, sizeName)
		} else {
			fieldname := field.Oneof.GoName
			if _, ok := oneofs[fieldname]; !ok {
				oneofs[fieldname] = struct{}{}
				p.P(`if vtmsg, ok := m.`, fieldname, `.(vtprotoMessage); ok {`)
				p.P(`n+=vtmsg.`, sizeName, `()`)
				p.P(`}`)
			}
		}
	}
	p.P(`if m.unknownFields != nil {`)
	p.P(`n+=len(m.unknownFields)`)
	p.P(`}`)
	p.P(`return n`)
	p.P(`}`)
	p.P()

	for _, field := range message.Fields {
		if field.Oneof == nil {
			continue
		}
		ccTypeName := field.GoIdent
		p.P(`func (m *`, ccTypeName, `) `, sizeName, `() (n int) {`)
		p.P(`if m == nil {`)
		p.P(`return 0`)
		p.P(`}`)
		p.P(`var l int`)
		p.P(`_ = l`)
		p.sizeForField(false, field, sizeName)
		p.P(`return n`)
		p.P(`}`)
	}
}

func (p *vtprotofile) SizeHelpers() {
	p.P(`
	func sov(x uint64) (n int) {
                return (`, p.Ident("math/bits", "Len64"), `(x | 1) + 6)/ 7
	}`)
	p.P(`func soz(x uint64) (n int) {
		return sov(uint64((x << 1) ^ uint64((int64(x) >> 63))))
	}`)
}
