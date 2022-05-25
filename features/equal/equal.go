// Copyright (c) 2022 PlanetScale Inc. All rights reserved.

package equal

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/planetscale/vtprotobuf/generator"
)

func init() {
	generator.RegisterFeature("equal", func(gen *generator.GeneratedFile) generator.FeatureGenerator {
		return &equal{GeneratedFile: gen}
	})
}

type equal struct {
	*generator.GeneratedFile
	once bool
}

var _ generator.FeatureGenerator = (*equal)(nil)

func (p *equal) Name() string     { return "equal" }
func (p *equal) GenerateHelpers() {}
func (p *equal) GenerateFile(file *protogen.File) bool {
	proto3 := file.Desc.Syntax() == protoreflect.Proto3
	for _, message := range file.Messages {
		p.message(proto3, message)
	}
	return p.once
}

const equalName = "EqualVT"

func (p *equal) message(proto3 bool, message *protogen.Message) {
	for _, nested := range message.Messages {
		p.message(proto3, nested)
	}

	if message.Desc.IsMapEntry() {
		return
	}

	p.once = true

	ccTypeName := message.GoIdent
	p.P(`func (this *`, ccTypeName, `) `, equalName, `(that *`, ccTypeName, `) bool {`)

	p.P(`if this == nil {`)
	p.P(`	return that == nil || that.String() == ""`)
	p.P(`} else if that == nil {`)
	p.P(`	return this.String() == ""`)
	p.P(`}`)

	sort.Slice(message.Fields, func(i, j int) bool {
		return message.Fields[i].Desc.Number() < message.Fields[j].Desc.Number()
	})

	{
		oneofs := make(map[string]struct{}, len(message.Fields))
		scoped := false
		for _, field := range message.Fields {
			oneof := field.Oneof != nil && !field.Oneof.Desc.IsSynthetic()
			nullable := field.Message != nil || (field.Oneof != nil && field.Oneof.Desc.IsSynthetic()) || (!proto3 && !oneof)
			if oneof {
				fieldname := field.Oneof.GoName
				if _, ok := oneofs[fieldname]; !ok {
					oneofs[fieldname] = struct{}{}
					p.P(`if this.`, fieldname, ` == nil && that.`, fieldname, ` != nil {`)
					p.P(`	return false`)
					p.P(`} else if this.`, fieldname, ` != nil {`)
					p.P(`	if that.`, fieldname, ` == nil {`)
					p.P(`		return false`)
					p.P(`	}`)
					scoped = true
				}

				p.oneof(field, nullable)
			}
		}
		if scoped {
			p.P(`}`)
		}
	}

	for _, field := range message.Fields {
		oneof := field.Oneof != nil && !field.Oneof.Desc.IsSynthetic()
		nullable := field.Message != nil || (field.Oneof != nil && field.Oneof.Desc.IsSynthetic()) || (!proto3 && !oneof)
		if !oneof {
			p.field(field, nullable)
		}
	}

	p.P(`return string(this.unknownFields) == string(that.unknownFields)`)
	p.P(`}`)
	p.P()
}

func (p *equal) oneof(field *protogen.Field, nullable bool) {
	fieldname := field.GoName

	getter := fmt.Sprintf("Get%s()", fieldname)
	kind := field.Desc.Kind()
	switch {
	case isScalar(kind):
		p.compareScalar(getter, nullable)
	case kind == protoreflect.BytesKind:
		p.compareBytes(getter)
	case kind == protoreflect.MessageKind || kind == protoreflect.GroupKind:
		goTyp, _ := p.FieldGoType(field)
		p.compareCall(getter, "", goTyp, field.Message)
	default:
		panic("not implemented")
	}
}

func (p *equal) field(field *protogen.Field, nullable bool) {
	fieldname := field.GoName

	repeated := field.Desc.Cardinality() == protoreflect.Repeated
	if repeated {
		p.P(`if len(this.`, fieldname, `) != len(that.`, fieldname, `) {`)
		p.P(`	return false`)
		p.P(`}`)
	}

	kind := field.Desc.Kind()
	switch {
	case isScalar(kind):
		if !repeated {
			p.compareScalar(fieldname, nullable)
			return
		}
		p.P(`for i := range this.`, fieldname, ` {`)
		p.compareScalar(fieldname+"[i]", false)
		p.P(`}`)

	case kind == protoreflect.BytesKind:
		if !repeated {
			p.compareBytes(fieldname)
			return
		}
		p.P(`for i := range this.`, fieldname, ` {`)
		p.compareBytes(fieldname + "[i]")
		p.P(`}`)

	case kind == protoreflect.MessageKind || kind == protoreflect.GroupKind:
		if !repeated {
			goTyp, _ := p.FieldGoType(field)
			p.compareCall(field.GoName, "", goTyp, field.Message)
			return
		}
		p.P(`for i := range this.`, fieldname, ` {`)
		if mv := field.Desc.MapValue(); mv != nil {
			kind := mv.Kind()
			switch {
			case isScalar(kind):
				p.compareScalar(fieldname+"[i]", false)
			case kind == protoreflect.BytesKind:
				p.compareBytes(fieldname + "[i]")
			case kind == protoreflect.MessageKind || kind == protoreflect.GroupKind:
				valueField := field.Message.Fields[1]
				goTypV, _ := p.FieldGoType(valueField)
				p.compareCall(field.GoName, "[i]", goTypV, valueField.Message)
			default:
				panic("not implemented")
			}
		} else {
			goTyp, _ := p.FieldGoType(field)
			if field.Desc.IsList() {
				goTyp = goTyp[len("[]"):]
			}
			p.compareCall(field.GoName, "[i]", goTyp, field.Message)
		}
		p.P(`}`)

	default:
		panic("not implemented")
	}
}

func (p *equal) compareScalar(fieldname string, nullable bool) {
	if nullable {
		p.P(`if p, q := this.`, fieldname, `, that.`, fieldname, `; (p == nil && q != nil) || (p != nil && (q == nil || *p != *q)) {`)
	} else {
		p.P(`if this.`, fieldname, ` != that.`, fieldname, ` {`)
	}
	p.P(`	return false`)
	p.P(`}`)
}

func (p *equal) compareBytes(fieldname string) {
	// Inlined call to bytes.Equal()
	p.P(`if string(this.`, fieldname, `) != string(that.`, fieldname, `) {`)
	p.P(`	return false`)
	p.P(`}`)
}

func (p *equal) compareCall(fieldname string, suffix, ccTypeName string, msg *protogen.Message) {
	if msg != nil && msg.Desc != nil && msg.Desc.ParentFile() != nil && p.IsLocalMessage(msg) {
		p.P(`if !this.`, fieldname, suffix, `.`, equalName, `(that.`, fieldname, suffix, `) {`)
		p.P(`	return false`)
		p.P(`}`)
		return
	}
	p.P(`if equal, ok := interface{}(this.`, fieldname, suffix, `).(interface { `, equalName, `(`, ccTypeName, `) bool }); ok {`)
	p.P(`	if !equal.`, equalName, `(that.`, fieldname, suffix, `) {`)
	p.P(`		return false`)
	p.P(`	}`)
	p.P(`} else if !`, p.Ident("google.golang.org/protobuf/proto", "Equal"), `(this.`, fieldname, suffix, `, that.`, fieldname, suffix, `) {`)
	p.P(`	return false`)
	p.P(`}`)
}

func isScalar(kind protoreflect.Kind) bool {
	switch kind {
	case
		protoreflect.BoolKind,
		protoreflect.StringKind,
		protoreflect.DoubleKind, protoreflect.Fixed64Kind, protoreflect.Sfixed64Kind,
		protoreflect.FloatKind, protoreflect.Fixed32Kind, protoreflect.Sfixed32Kind,
		protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Sint64Kind,
		protoreflect.Int32Kind, protoreflect.Uint32Kind, protoreflect.Sint32Kind,
		protoreflect.EnumKind:
		return true
	}
	return false
}
