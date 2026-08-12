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
	ccTypeName := message.GoIdent

	p.P(`var vtprotoPool_`, ccTypeName, ` = `, p.Ident("sync", "Pool"), `{`)
	p.P(`New: func() interface{} {`)
	p.P(`return &`, ccTypeName, `{}`)
	p.P(`},`)
	p.P(`}`)

	p.P(`func (m *`, ccTypeName, `) ResetVT() {`)
	p.P(`if m != nil {`)
	var saved []*protogen.Field
	for _, field := range message.Fields {
		fieldName := p.FieldName(field)

		if field.Desc.IsList() && p.FieldStorageIsPointer(field) {
			// Keep the pointer so the slice survives m.Reset().
			ptr := fmt.Sprintf("f%d", len(saved))
			p.P(ptr, ` := m.`, fieldName)
			p.P(`if `, ptr, ` != nil {`)
			p.P(`for _, mm := range *`, ptr, `{`)
			if p.ShouldPool(field.Message) {
				p.P(`mm.ResetVT()`)
			} else {
				p.P(`mm.Reset()`)
			}
			p.P(`}`)
			p.P(`*`, ptr, ` = (*`, ptr, `)[:0]`)
			p.P(`}`)
			saved = append(saved, field)
		} else if field.Desc.IsList() {
			switch field.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
				p.P(`for _, mm := range m.`, fieldName, `{`)
				if p.ShouldPool(field.Message) {
					p.P(`mm.ResetVT()`)
				} else {
					p.P(`mm.Reset()`)
				}
				p.P(`}`)
			case protoreflect.BytesKind, protoreflect.StringKind:
				p.P(`clear(m.`, fieldName, `)`)
			}
			p.P(fmt.Sprintf("f%d", len(saved)), ` := m.`, fieldName, `[:0]`)
			saved = append(saved, field)
		} else if field.Oneof != nil && !field.Oneof.Desc.IsSynthetic() {
			if p.ShouldPool(field.Message) {
				p.P(`if oneof, ok := `, p.OneofStore("m", field.Oneof), `.(*`, p.OneofWrapperIdent(field), `); ok {`)
				p.P(`oneof.`, field.GoName, `.ReturnToVTPool()`)
				p.P(`}`)
			}
		} else {
			switch field.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
				if !field.Desc.IsMap() && p.ShouldPool(field.Message) {
					p.P(`m.`, fieldName, `.ReturnToVTPool()`)
				}
			case protoreflect.BytesKind:
				// Retaining the buffer would contradict the bit m.Reset() cleared.
				if p.UsesPresenceBit(field) {
					break
				}
				p.P(fmt.Sprintf("f%d", len(saved)), ` := m.`, fieldName, `[:0]`)
				saved = append(saved, field)
			}
		}
	}

	p.P(`m.Reset()`)
	for i, field := range saved {
		p.P(`m.`, p.FieldName(field), ` = `, fmt.Sprintf("f%d", i))
	}
	p.P(`}`)
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
