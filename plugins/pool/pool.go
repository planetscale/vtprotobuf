package pool

import (
	"fmt"
	"vitess.io/vtprotobuf/plugins/common"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func init() {
	common.RegisterPlugin(func(gen *common.VTGeneratedFile) common.Plugin {
		return &pooler{VTGeneratedFile: gen}
	})
}

type pooler struct {
	*common.VTGeneratedFile
	once bool
}

var _ common.Plugin = (*pooler)(nil)

func (p *pooler) Name() string {
	return "pool"
}

func (p *pooler) GenerateHelpers() {
}
func (p *pooler) GenerateFile(file *protogen.File) bool {
	for _, message := range file.Messages {
		p.message(message)
	}
	return p.once
}

func (p *pooler) message(message *protogen.Message) {
	for _, nested := range message.Messages {
		p.message(nested)
	}

	if message.Desc.IsMapEntry() || !p.ShouldPool(message) {
		return
	}

	p.once = true
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
				if p.ShouldPool(field.Message) {
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
				if p.ShouldPool(field.Message) {
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
