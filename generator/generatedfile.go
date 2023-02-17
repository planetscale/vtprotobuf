// Copyright (c) 2021 PlanetScale Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

import (
	"fmt"

	"github.com/planetscale/vtprotobuf/vtproto"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type GeneratedFile struct {
	*protogen.GeneratedFile
	Ext           *Extensions
	LocalPackages map[string]bool

	helpers map[string]bool
}

func (p *GeneratedFile) Helper(name string, generate func(p *GeneratedFile)) {
	if p.helpers[name] {
		return
	}
	generate(p)
	p.helpers[name] = true
}

func (p *GeneratedFile) Ident(path, ident string) string {
	return p.QualifiedGoIdent(protogen.GoImportPath(path).Ident(ident))
}

func (p *GeneratedFile) ShouldPool(message *protogen.Message) bool {
	if message == nil {
		return false
	}
	if p.Ext.Poolable[message.GoIdent] {
		return true
	}
	ext := proto.GetExtension(message.Desc.Options(), vtproto.E_Mempool)
	if mempool, ok := ext.(bool); ok {
		return mempool
	}
	return false
}

func (p *GeneratedFile) Alloc(vname string, message *protogen.Message) {
	if p.ShouldPool(message) {
		p.P(vname, " := ", message.GoIdent, `FromVTPool()`)
	} else {
		p.P(vname, " := new(", message.GoIdent, `)`)
	}
}

func (p *GeneratedFile) FieldGoType(field *protogen.Field) (goType string, pointer bool) {
	if field.Desc.IsWeak() {
		return "struct{}", false
	}

	pointer = field.Desc.HasPresence()
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		goType = "bool"
	case protoreflect.EnumKind:
		goType = p.QualifiedGoIdent(field.Enum.GoIdent)
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		goType = "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		goType = "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		goType = "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		goType = "uint64"
	case protoreflect.FloatKind:
		goType = "float32"
	case protoreflect.DoubleKind:
		goType = "float64"
	case protoreflect.StringKind:
		goType = "string"
	case protoreflect.BytesKind:
		goType = "[]byte"
		pointer = false // rely on nullability of slices for presence
	case protoreflect.MessageKind, protoreflect.GroupKind:
		goType = "*" + p.QualifiedGoIdent(field.Message.GoIdent)
		pointer = false // pointer captured as part of the type
	}
	switch {
	case field.Desc.IsList():
		return "[]" + goType, false
	case field.Desc.IsMap():
		keyType, _ := p.FieldGoType(field.Message.Fields[0])
		valType, _ := p.FieldGoType(field.Message.Fields[1])
		return fmt.Sprintf("map[%v]%v", keyType, valType), false
	}
	return goType, pointer
}

func (p *GeneratedFile) IsLocalMessage(message *protogen.Message) bool {
	pkg := string(message.Desc.ParentFile().Package())
	return p.LocalPackages[pkg]
}

func (p *GeneratedFile) IsWellKnownMessage(message *protogen.Message) bool {
	wellknown := map[string]bool{"google.protobuf.Timestamp": true}
	return wellknown[p.MessageID(message)]
}

func (p *GeneratedFile) MessageID(message *protogen.Message) string {
	return fmt.Sprintf("%s.%s", message.Desc.ParentFile().Package(), message.Desc.Name())
}
