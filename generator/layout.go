// Copyright (c) 2021 PlanetScale Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generator

import (
	"strconv"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/gofeaturespb"
)

var protoimplPackage = protogen.GoImportPath("google.golang.org/protobuf/runtime/protoimpl")

// The helpers here are the only place that knows how a message's Go struct is
// laid out. HYBRID counts as OPEN: the variant protoc-gen-go emits by default
// for a hybrid file has the open layout, and the file generated alongside it
// carries a !protoopaque build constraint.

func isOpaque(message *protogen.Message) bool {
	return message != nil && message.APILevel == gofeaturespb.GoFeatures_API_OPAQUE
}

// isLazy reports whether the field is marked for lazy decoding. protoreflect
// does not expose this, so assert for it the way protoc-gen-go does.
func isLazy(field *protogen.Field) bool {
	lazy, ok := field.Desc.(interface{ IsLazy() bool })
	return ok && lazy.IsLazy()
}

// FieldName returns the name of the Go struct field backing field. Features
// build accesses as "m." + FieldName(field), which is what keeps the opaque
// name mangling confined to this file.
func (p *GeneratedFile) FieldName(field *protogen.Field) string {
	if isOpaque(field.Parent) {
		return "xxx_hidden_" + field.GoName
	}
	return field.GoName
}

func (p *GeneratedFile) OneofName(oneof *protogen.Oneof) string {
	if isOpaque(oneof.Parent) {
		return "xxx_hidden_" + oneof.GoName
	}
	return oneof.GoName
}

func (p *GeneratedFile) fieldStore(recv string, field *protogen.Field) string {
	return recv + "." + p.FieldName(field)
}

// FieldStorageIsPointer reports whether the Go struct field is a pointer.
// Pointers that are part of the type itself, as in messages and slices, do not
// count.
func (p *GeneratedFile) FieldStorageIsPointer(field *protogen.Field) bool {
	if !isOpaque(field.Parent) {
		// Mirrors FieldGoType's pointer computation without calling it: FieldGoType
		// resolves QualifiedGoIdent for message/enum types as a side effect, which
		// would register an unwanted import here even when the caller only wants
		// the bool and never prints the resolved type.
		switch {
		case field.Desc.IsWeak():
			return false
		case field.Desc.IsList(), field.Desc.IsMap():
			return false
		case field.Desc.Kind() == protoreflect.BytesKind:
			return false
		case field.Desc.Kind() == protoreflect.MessageKind, field.Desc.Kind() == protoreflect.GroupKind:
			return false
		default:
			return field.Desc.HasPresence()
		}
	}
	switch {
	case field.Desc.IsMap():
		return false
	case field.Desc.IsList():
		return field.Message != nil
	case field.Desc.Kind() == protoreflect.StringKind:
		// Strings keep their pointer even though a presence bit also tracks them.
		return field.Desc.HasPresence()
	default:
		return false
	}
}

// FieldSliceExpr returns the expression to read a repeated field as a plain
// slice, emitting a local first if the layout stores the slice behind a pointer.
func (p *GeneratedFile) FieldSliceExpr(recv string, field *protogen.Field) string {
	store := p.fieldStore(recv, field)
	if !field.Desc.IsList() || field.Message == nil || field.Desc.IsMap() || !isOpaque(field.Parent) {
		return store
	}
	local := "vt" + recv + field.GoName
	goType, _ := p.FieldGoType(field)
	p.P(`var `, local, ` `, goType)
	p.P(`if `, local, `p := `, store, `; `, local, `p != nil {`)
	p.P(local, ` = *`, local, `p`)
	p.P(`}`)
	return local
}

func (p *GeneratedFile) presenceWord(recv string, field *protogen.Field) (word, index string) {
	idx := presenceIndex(field)
	return "&(" + recv + ".XXX_presence[" + strconv.Itoa(idx/32) + "])", strconv.Itoa(idx)
}

func (p *GeneratedFile) FieldPresent(recv string, field *protogen.Field) string {
	if p.UsesPresenceBit(field) {
		word, idx := p.presenceWord(recv, field)
		return p.QualifiedGoIdent(protoimplPackage.Ident("X")) + ".Present(" + word + ", " + idx + ")"
	}
	return p.fieldStore(recv, field) + " != nil"
}

func (p *GeneratedFile) FieldAbsent(recv string, field *protogen.Field) string {
	if p.UsesPresenceBit(field) {
		return "!" + p.FieldPresent(recv, field)
	}
	return p.fieldStore(recv, field) + " == nil"
}

// FieldSetPresent emits the statement marking field as present. SetPresent
// rather than SetPresentNonAtomic: UnmarshalVT may run on a message another
// goroutine holds, and neighbouring fields share a presence word.
func (p *GeneratedFile) FieldSetPresent(recv string, field *protogen.Field) {
	if !p.UsesPresenceBit(field) {
		return
	}
	word, idx := p.presenceWord(recv, field)
	p.P(protoimplPackage.Ident("X"), ".SetPresent(", word, ", ", idx, ", ",
		strconv.Itoa(numPresenceFields(field.Parent)), ")")
}

func (p *GeneratedFile) CopyPresence(dst, src string, message *protogen.Message) {
	if !p.hasPresenceArray(message) {
		return
	}
	p.P(dst, ".XXX_presence = ", src, ".XXX_presence")
}

func (p *GeneratedFile) hasPresenceArray(message *protogen.Message) bool {
	for _, field := range message.Fields {
		if p.UsesPresenceBit(field) {
			return true
		}
	}
	return false
}

func (p *GeneratedFile) OneofStore(recv string, oneof *protogen.Oneof) string {
	return recv + "." + p.OneofName(oneof)
}

// OneofWrapperIdent returns the wrapper struct type for a oneof member, which
// the opaque API unexports.
func (p *GeneratedFile) OneofWrapperIdent(field *protogen.Field) protogen.GoIdent {
	id := field.GoIdent
	if isOpaque(field.Parent) {
		r, sz := utf8.DecodeRuneInString(id.GoName)
		id.GoName = string(unicode.ToLower(r)) + id.GoName[sz:]
	}
	return id
}

// UsesPresenceBit reports whether presence lives in the XXX_presence bitfield
// rather than in a nil pointer. Mirrors filedesc.UsePresenceForField, which is
// internal.
func (p *GeneratedFile) UsesPresenceBit(field *protogen.Field) bool {
	if !isOpaque(field.Parent) {
		return false
	}
	switch {
	case field.Desc.IsList(), field.Desc.IsMap():
		return false
	case field.Oneof != nil && !field.Oneof.Desc.IsSynthetic():
		return false
	case field.Message != nil:
		return false
	}
	return field.Desc.HasPresence()
}

// presenceIndex returns the field's slot in the bitfield: its position among the
// emitted struct fields, where a whole oneof takes one slot accounted at its
// last member. Mirrors opaqueFieldPresenceIndex in protoc-gen-go.
func presenceIndex(field *protogen.Field) int {
	idx := 0
	for _, f := range field.Parent.Fields {
		if f == field {
			break
		}
		if f.Oneof == nil || isLastOneofField(f) {
			idx++
		}
	}
	return idx
}

func numPresenceFields(message *protogen.Message) int {
	if len(message.Fields) == 0 {
		return 0
	}
	return presenceIndex(message.Fields[len(message.Fields)-1]) + 1
}

func isLastOneofField(field *protogen.Field) bool {
	fields := field.Oneof.Fields
	return fields[len(fields)-1] == field
}
