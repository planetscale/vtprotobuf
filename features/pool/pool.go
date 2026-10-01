package pool

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/planetscale/vtprotobuf/generator"
)

func init() {
	generator.RegisterFeature("pool", func(gen *generator.GeneratedFile) generator.FeatureGenerator {
		return &pool{GeneratedFile: gen}
	})
}

type pool struct {
	*generator.GeneratedFile
	once bool
}

var _ generator.FeatureGenerator = (*pool)(nil)

func (p *pool) GenerateFile(file *protogen.File) bool {
	for _, message := range file.Messages {
		p.message(message)
	}
	return p.once
}

func (p *pool) message(message *protogen.Message) {
	for _, nested := range message.Messages {
		p.message(nested)
	}

	if message.Desc.IsMapEntry() || !p.ShouldPool(message) {
		return
	}

	p.once = true
	ccTypeName := p.PoolMessageIdent(message)

	p.P(`var vtprotoPool_`, ccTypeName.GoName, ` = `, p.Ident("sync", "Pool"), `{`)
	p.P(`New: func() interface{} {`)
	p.P(`return &`, ccTypeName, `{}`)
	p.P(`},`)
	p.P(`}`)

	p.P(`func (m *`, ccTypeName, `) ResetVT() {`)
	p.P(`if m != nil {`)
	var saved []*protogen.Field
	for _, field := range message.Fields {
		fieldName := field.GoName

		if field.Desc.IsList() {
			switch field.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
				p.P(`for _, mm := range m.`, fieldName, `{`)
				if p.ShouldPool(field.Message) {
					p.PResetVTElement("mm", field.Message)
				} else {
					p.P(`mm.Reset()`)
				}
				p.P(`}`)
			case protoreflect.BytesKind, protoreflect.StringKind:
				p.P(`clear(m.`, fieldName, `)`)
			}
			p.P(fmt.Sprintf("f%d", len(saved)), ` := m.`, fieldName, `[:0]`)
			saved = append(saved, field)
		} else if field.Desc.IsMap() {
			// Pooled message values go back to their pool; other values are dropped
			// by clear. The map itself is kept so its buckets are reused.
			if valueField := field.Message.Fields[1]; valueField.Message != nil && p.ShouldPool(valueField.Message) {
				p.P(`for _, v := range m.`, fieldName, `{`)
				p.P(`if v != nil {`)
				p.PReturnToVTPool("v", valueField.Message)
				p.P(`}`)
				p.P(`}`)
			}
			p.P(`clear(m.`, fieldName, `)`)
			p.P(fmt.Sprintf("f%d", len(saved)), ` := m.`, fieldName)
			saved = append(saved, field)
		} else if field.Oneof != nil && !field.Oneof.Desc.IsSynthetic() {
			if p.ShouldPool(field.Message) {
				p.P(`if oneof, ok := m.`, field.Oneof.GoName, `.(*`, field.GoIdent, `); ok {`)
				p.PReturnToVTPool("oneof."+fieldName, field.Message)
				p.P(`}`)
			}
		} else {
			switch field.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
				if p.ShouldPool(field.Message) {
					p.PReturnToVTPool("m."+fieldName, field.Message)
				}
			case protoreflect.BytesKind:
				p.P(fmt.Sprintf("f%d", len(saved)), ` := m.`, fieldName, `[:0]`)
				saved = append(saved, field)
			}
		}
	}

	if p.IsLocalWrapper(message) {
		p.P(`(*`, p.QualifiedGoIdent(message.GoIdent), `)(m).Reset()`)
	} else {
		p.P(`m.Reset()`)
	}
	for i, field := range saved {
		p.P(`m.`, field.GoName, ` = `, fmt.Sprintf("f%d", i))
	}
	p.P(`}`)
	p.P(`}`)

	p.P(`func (m *`, ccTypeName, `) ReturnToVTPool() {`)
	p.P(`if m != nil {`)
	p.P(`m.ResetVT()`)
	p.P(`vtprotoPool_`, ccTypeName.GoName, `.Put(m)`)
	p.P(`}`)
	p.P(`}`)

	p.P(`func `, ccTypeName.GoName, `FromVTPool() *`, ccTypeName, `{`)
	p.P(`return vtprotoPool_`, ccTypeName.GoName, `.Get().(*`, ccTypeName, `)`)
	p.P(`}`)
}
