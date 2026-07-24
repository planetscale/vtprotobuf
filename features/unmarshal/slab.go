// Copyright (c) 2021 PlanetScale Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unmarshal

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/planetscale/vtprotobuf/generator"
)

func init() {
	generator.RegisterFeature("unmarshal_slab", func(gen *generator.GeneratedFile) generator.FeatureGenerator {
		return &unmarshal{GeneratedFile: gen, slab: true}
	})
}

// The unmarshal_slab feature emits, for every message that opts in with
// `option (vtproto.slab) = true;` (or a `slab=` plugin option), an
// UnmarshalVTSlab entry point that decodes exactly like UnmarshalVT but
// allocates nested message elements and scalar presence pointers from
// chunked slabs, so N element allocations collapse into O(log N) chunk
// allocations. A counting pre-pass over the top-level wire format reserves
// exactly-sized chunks for the repeated message fields before decoding, and
// payloads below protohelpers.SlabUnmarshalThreshold skip the arena and fall
// back to UnmarshalVT (for tiny messages the setup outweighs the savings).
//
// Slab elements are handed out once and never reused, so decoded messages
// behave like independently heap-allocated ones — except that retaining any
// pointer into the message graph pins that pointer's whole chunk. The entry
// point is therefore generated only for messages that explicitly opt in, and
// is meant for messages whose lifetime ends with the RPC that carried them.
//
// The transform covers the message graph reachable from an opted-in root
// through message-typed fields declared in the same .proto file; fields of
// types from other files, well-known types, and groups keep the stock
// decoding path (including unknown-field retention, which is preserved
// everywhere).

// generateSlabFile drives the unmarshal_slab feature for one file: it
// computes the set of messages reachable from the opted-in roots, collects
// the slab element types they need into a per-file arena struct, and then
// reuses the stock unmarshal emission (with allocation sites redirected to
// the arena) to generate the arena-threaded decode methods.
func (p *unmarshal) generateSlabFile(file *protogen.File, proto3 bool) bool {
	if p.Wrapper() {
		return false
	}

	var all []*protogen.Message
	var collect func(messages []*protogen.Message)
	collect = func(messages []*protogen.Message) {
		for _, m := range messages {
			if !m.Desc.IsMapEntry() {
				all = append(all, m)
			}
			collect(m.Messages)
		}
	}
	collect(file.Messages)

	var roots []*protogen.Message
	for _, m := range all {
		if p.ShouldSlab(m) {
			if p.ShouldPool(m) {
				panic(fmt.Sprintf("message %s enables both memory pooling and slab unmarshalling; the two cannot be combined", m.Desc.FullName()))
			}
			roots = append(roots, m)
		}
	}
	if len(roots) == 0 {
		return false
	}

	p.slabFile = file
	p.slabClosure = make(map[*protogen.Message]bool)
	queue := append([]*protogen.Message(nil), roots...)
	for _, m := range roots {
		p.slabClosure[m] = true
	}
	for len(queue) > 0 {
		m := queue[0]
		queue = queue[1:]
		for _, field := range m.Fields {
			if target := p.slabFieldTarget(field); target != nil && !p.slabClosure[target] {
				p.slabClosure[target] = true
				queue = append(queue, target)
			}
		}
	}

	// Register the arena slab fields in declaration order so the generated
	// output is deterministic.
	p.slabFieldByType = make(map[string]string)
	p.slabFieldNames = make(map[string]bool)
	for _, m := range all {
		if !p.slabClosure[m] {
			continue
		}
		for _, field := range m.Fields {
			if target := p.slabFieldTarget(field); target != nil {
				p.registerSlabField(target.GoIdent.GoName)
			}
			if typ, ok := p.pointerScalarType(proto3, field); ok {
				p.registerSlabField(typ)
			}
		}
	}

	p.slabArena = "slabArena_" + strings.TrimPrefix(file.GoDescriptorIdent.GoName, "File_")

	p.P(`// `, p.slabArena, ` holds the slabs backing the UnmarshalVTSlab entry`)
	p.P(`// points of this file. One arena lives for the duration of a single`)
	p.P(`// top-level unmarshal call.`)
	p.P(`type `, p.slabArena, ` struct {`)
	for _, typ := range p.slabFieldOrder {
		p.P(p.slabFieldByType[typ], ` `, p.Helper("Slab"), `[`, typ, `]`)
	}
	p.P(`}`)
	p.P()

	for _, message := range file.Messages {
		p.message(proto3, message)
	}
	return p.once
}

// slabFieldTarget returns the message type that unmarshalling this field
// allocates from a slab (the element type for repeated fields, the value
// type for maps, the field type otherwise), or nil if the field keeps the
// stock allocation path. Only message types declared in the same file
// participate: their arena-threaded decode methods are generated alongside
// the arena type itself, so the output stays self-contained.
func (p *unmarshal) slabFieldTarget(field *protogen.Field) *protogen.Message {
	target := field.Message
	if field.Desc.IsMap() {
		vf := field.Message.Fields[1]
		if vf.Desc.Kind() != protoreflect.MessageKind {
			return nil
		}
		target = vf.Message
	} else if field.Desc.Kind() != protoreflect.MessageKind {
		return nil
	}
	if target.Desc.ParentFile().Path() != p.slabFile.Desc.Path() {
		return nil
	}
	return target
}

// pointerScalarType reports whether unmarshalling this field stores a
// pointer to a scalar (a proto2 or explicit-presence scalar, enum or string
// field), and returns the pointed-to Go type. These are the `m.F = &v` sites
// in the stock decoder that slab mode redirects to per-type slabs.
func (p *unmarshal) pointerScalarType(proto3 bool, field *protogen.Field) (string, bool) {
	if field.Desc.IsMap() || field.Desc.IsList() {
		return "", false
	}
	if field.Oneof != nil && !field.Oneof.Desc.IsSynthetic() {
		return "", false
	}
	nullable := field.Oneof != nil && field.Oneof.Desc.IsSynthetic()
	if proto3 && !nullable {
		return "", false
	}
	switch field.Desc.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind, protoreflect.BytesKind:
		return "", false
	}
	return p.noStarOrSliceType(field), true
}

// registerSlabField assigns (or returns) the arena field name for a slab
// element type, keeping names unique even if two distinct type expressions
// sanitize to the same identifier.
func (p *unmarshal) registerSlabField(typ string) string {
	if name, ok := p.slabFieldByType[typ]; ok {
		return name
	}
	base := "f_" + sanitizeSlabIdent(typ)
	name := base
	for i := 2; p.slabFieldNames[name]; i++ {
		name = base + strconv.Itoa(i)
	}
	p.slabFieldNames[name] = true
	p.slabFieldByType[typ] = name
	p.slabFieldOrder = append(p.slabFieldOrder, typ)
	return name
}

// slabField returns the arena field name registered for a slab element type.
func (p *unmarshal) slabField(typ string) string {
	name, ok := p.slabFieldByType[typ]
	if !ok {
		panic(fmt.Sprintf("no slab arena field registered for type %q", typ))
	}
	return name
}

func sanitizeSlabIdent(typ string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			return r
		}
		return '_'
	}, typ)
}

// storePtrVar stores a pointer to the decoded scalar variable into the
// field: the stock decoder takes the address of the stack temporary, slab
// mode copies the value into a per-type slab instead.
func (p *unmarshal) storePtrVar(fieldname, typ, varName string) {
	if p.slab {
		p.P(`m.`, fieldname, ` = a.`, p.slabField(typ), `.NextValue(`, varName, `)`)
	} else {
		p.P(`m.`, fieldname, ` = &`, varName)
	}
}

// storePtrExpr is storePtrVar for sites where the stock decoder materializes
// the value through a named temporary (`tmpName := expr; m.F = &tmpName`).
func (p *unmarshal) storePtrExpr(fieldname, typ, tmpName, expr string) {
	if p.slab {
		p.P(`m.`, fieldname, ` = a.`, p.slabField(typ), `.NextValue(`, expr, `)`)
	} else {
		p.P(tmpName, ` := `, expr)
		p.P(`m.`, fieldname, ` = &`, tmpName)
	}
}

// slabEntryPoint emits the public UnmarshalVTSlab method for an opted-in
// message: the small-payload bypass, the counting pre-pass over the
// top-level wire format, the slab reservations derived from those counts,
// and the call into the arena-threaded decoder.
func (p *unmarshal) slabEntryPoint(message *protogen.Message) {
	ccTypeName := message.GoIdent.GoName

	// The counted fields: top-level repeated message fields and maps with
	// message values, whose elements come from slabs. One linear scan
	// counts them all so their slabs can be reserved exactly-sized.
	type countedField struct {
		number int32
		target *protogen.Message
	}
	var counted []countedField
	for _, field := range message.Fields {
		if !field.Desc.IsMap() && !field.Desc.IsList() {
			continue
		}
		if target := p.slabFieldTarget(field); target != nil {
			counted = append(counted, countedField{int32(field.Desc.Number()), target})
		}
	}

	// Reservations: the counted types get their exact counts, and every type
	// reachable through chains of plain singular message fields gets the
	// originating count multiplied by how many instances one parent can hold
	// (an upper bound: singular fields may be absent). The root's own
	// singular chains contribute a constant. Anything else — nested repeated
	// fields, oneof variants, scalar slabs — falls back to chunked growth.
	type resRow struct {
		constant int
		coeff    []int
	}
	rows := make(map[*protogen.Message]*resRow)
	var resOrder []*protogen.Message
	row := func(t *protogen.Message) *resRow {
		r, ok := rows[t]
		if !ok {
			r = &resRow{coeff: make([]int, len(counted))}
			rows[t] = r
			resOrder = append(resOrder, t)
		}
		return r
	}

	memo := make(map[*protogen.Message]map[*protogen.Message]int)
	for _, field := range message.Fields {
		if target := p.singularFieldTarget(field); target != nil {
			row(target).constant++
			for t, n := range p.singularReach(target, memo, map[*protogen.Message]bool{}) {
				row(t).constant += n
			}
		}
	}
	for i, cf := range counted {
		row(cf.target).coeff[i]++
		for t, n := range p.singularReach(cf.target, memo, map[*protogen.Message]bool{}) {
			row(t).coeff[i] += n
		}
	}

	p.P(`// UnmarshalVTSlab is like UnmarshalVT, but allocates nested message`)
	p.P(`// elements and scalar presence pointers from chunked slabs sized by a`)
	p.P(`// counting pre-pass over the wire format. Retaining any pointer into`)
	p.P(`// the decoded message pins that pointer's whole chunk in memory: use`)
	p.P(`// it for messages whose lifetime ends with the RPC that carried them.`)
	p.P(`func (m *`, ccTypeName, `) UnmarshalVTSlab(dAtA []byte) error {`)
	p.P(`if len(dAtA) < `, p.Helper("SlabUnmarshalThreshold"), ` {`)
	p.P(`return m.UnmarshalVT(dAtA)`)
	p.P(`}`)
	p.P(`var a `, p.slabArena)
	if len(counted) > 0 {
		k := strconv.Itoa(len(counted))
		nums := make([]string, len(counted))
		sum := make([]string, len(counted))
		for i, cf := range counted {
			nums[i] = strconv.Itoa(int(cf.number))
			sum[i] = `counts[` + strconv.Itoa(i) + `]`
		}
		p.P(`fieldNums := [`, k, `]int32{`, strings.Join(nums, ", "), `}`)
		p.P(`var counts [`, k, `]int`)
		p.P(p.Helper("CountFields"), `(dAtA, fieldNums[:], counts[:])`)
		p.P(`if `, strings.Join(sum, `+`), ` < `, p.Helper("SlabUnmarshalMinCount"), ` {`)
		p.P(`return m.UnmarshalVT(dAtA)`)
		p.P(`}`)
	}
	for _, t := range resOrder {
		r := rows[t]
		var terms []string
		if r.constant > 0 {
			terms = append(terms, strconv.Itoa(r.constant))
		}
		for i, c := range r.coeff {
			switch {
			case c == 1:
				terms = append(terms, `counts[`+strconv.Itoa(i)+`]`)
			case c > 1:
				terms = append(terms, strconv.Itoa(c)+`*counts[`+strconv.Itoa(i)+`]`)
			}
		}
		p.P(`a.`, p.slabField(t.GoIdent.GoName), `.Reserve(`, strings.Join(terms, ` + `), `)`)
	}
	p.P(`return m.unmarshalVTSlab(dAtA, &a)`)
	p.P(`}`)
	p.P()
}

// singularFieldTarget returns the slab closure type of a plain singular
// message field (not repeated, not a map, not a member of a real oneof — at
// most one oneof variant is set, so reserving for all of them would
// systematically over-allocate), or nil.
func (p *unmarshal) singularFieldTarget(field *protogen.Field) *protogen.Message {
	if field.Desc.IsMap() || field.Desc.IsList() {
		return nil
	}
	if field.Oneof != nil && !field.Oneof.Desc.IsSynthetic() {
		return nil
	}
	if field.Desc.Kind() != protoreflect.MessageKind {
		return nil
	}
	return p.slabFieldTarget(field)
}

// singularReach returns, for every slab closure type, how many instances a
// single instance of m can allocate through chains of plain singular message
// fields — the multiplier used to turn a parent count into a reservation
// upper bound. Results are memoized per file; recursion cycles are cut by
// treating the back edge as zero (an under-estimate, which merely falls back
// to chunked growth).
func (p *unmarshal) singularReach(m *protogen.Message, memo map[*protogen.Message]map[*protogen.Message]int, inProgress map[*protogen.Message]bool) map[*protogen.Message]int {
	if r, ok := memo[m]; ok {
		return r
	}
	if inProgress[m] {
		return nil
	}
	inProgress[m] = true
	r := make(map[*protogen.Message]int)
	for _, field := range m.Fields {
		target := p.singularFieldTarget(field)
		if target == nil {
			continue
		}
		r[target]++
		for t, n := range p.singularReach(target, memo, inProgress) {
			r[t] += n
		}
	}
	delete(inProgress, m)
	memo[m] = r
	return r
}
