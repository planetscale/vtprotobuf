package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var poolable = map[protogen.GoIdent]bool{
	protogen.GoIdent{
		GoName:       "Row",
		GoImportPath: "vitess.io/vitess/go/vt/proto/query",
	}: true,
	protogen.GoIdent{
		GoName:       "VStreamRowsResponse",
		GoImportPath: "vitess.io/vitess/go/vt/proto/binlogdata",
	}: true,
}

func (p *vtprotofile) shouldPool(message *protogen.Message) bool {
	if message == nil {
		return false
	}
	return poolable[message.GoIdent]
}

func (p *vtprotofile) PoolMessage(message *protogen.Message) {
	if !p.shouldPool(message) {
		return
	}

	ccTypeName := message.GoIdent

	p.P(`var vtprotoPool_`, ccTypeName, ` = `, p.Ident("sync", "Pool"), `{`)
	p.P(`New: func() interface{} {`)
	p.P(`return &`, message.GoIdent, `{}`)
	p.P(`},`)
	p.P(`}`)

	p.P(`func (m *`, ccTypeName, `) ReturnToVTPool() {`)
	p.P(`if m != nil {`)

	var saved []*protogen.Field
	for _, field := range message.Fields {
		fieldName := field.GoName

		if field.Desc.IsList() {
			switch field.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
				if p.shouldPool(field.Message) {
					p.P(`for _, mm := range m.`, fieldName, `{`)
					p.P(`mm.ReturnToVTPool()`)
					p.P(`}`)
				}
			}
			p.P(fmt.Sprintf("f%d", len(saved)), ` := m.`, fieldName, `[:0]`)
			saved = append(saved, field)
		} else {
			switch field.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
				if p.shouldPool(field.Message) {
					p.P(`m.`, fieldName, `.ReturnToVTPool()`)
				}
			case protoreflect.BytesKind:
				p.P(fmt.Sprintf("f%d", len(saved)), ` := m.`, fieldName, `[:0]`)
				saved = append(saved, field)
			}
		}
	}

	p.P(`m.Reset()`)

	for i, field := range saved {
		p.P(`m.`, field.GoName, ` = `, fmt.Sprintf("f%d", i))
	}

	p.P(`vtprotoPool_`, ccTypeName, `.Put(m)`)
	p.P(`}`)
	p.P(`}`)

	p.P(`func `, ccTypeName, `FromVTPool() *`, ccTypeName, `{`)
	p.P(`return vtprotoPool_`, ccTypeName, `.Get().(*`, ccTypeName, `)`)
	p.P(`}`)
}
