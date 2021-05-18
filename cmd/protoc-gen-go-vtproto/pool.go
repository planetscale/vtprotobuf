package main

import (
	"fmt"

	"vitess.io/vtprotobuf/vtproto"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func (p *vtprotofile) shouldPool(message *protogen.Message) bool {
	if message == nil {
		return false
	}
	if p.mempool[message.GoIdent] {
		return true
	}
	ext := proto.GetExtension(message.Desc.Options(), vtproto.E_Mempool)
	if mempool, ok := ext.(bool); ok {
		return mempool
	}
	return false
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

	p.P(`func (m *`, ccTypeName, `) ResetVT() {`)
	var saved []*protogen.Field
	for _, field := range message.Fields {
		fieldName := field.GoName

		if field.Desc.IsList() {
			switch field.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
				if p.shouldPool(field.Message) {
					p.P(`for _, mm := range m.`, fieldName, `{`)
					p.P(`mm.ResetVT()`)
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
	p.P(`}`)

	p.P(`func (m *`, ccTypeName, `) ReturnToVTPool() {`)
	p.P(`if m != nil {`)
	p.P(`m.ResetVT()`)
	p.P(`vtprotoPool_`, ccTypeName, `.Put(m)`)
	p.P(`}`)
	p.P(`}`)

	p.P(`func `, ccTypeName, `FromVTPool() *`, ccTypeName, `{`)
	p.P(`return vtprotoPool_`, ccTypeName, `.Get().(*`, ccTypeName, `)`)
	p.P(`}`)
}
