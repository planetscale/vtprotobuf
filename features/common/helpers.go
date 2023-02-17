// Commercial secret, LLC "RevTech". Refer to CONFIDENTIAL file in the root for details

package common

import (
	"github.com/planetscale/vtprotobuf/generator"
)

func HelperEncodeVarint(p *generator.GeneratedFile) {
	p.Helper("encodeVarint", func(p *generator.GeneratedFile) {
		p.P(`func encodeVarint(dAtA []byte, offset int, v uint64) int {`)
		p.P(`offset -= sov(v)`)
		p.P(`base := offset`)
		p.P(`for v >= 1<<7 {`)
		p.P(`dAtA[offset] = uint8(v&0x7f|0x80)`)
		p.P(`v >>= 7`)
		p.P(`offset++`)
		p.P(`}`)
		p.P(`dAtA[offset] = uint8(v)`)
		p.P(`return base`)
		p.P(`}`)
	})
}

func HelperSOV(p *generator.GeneratedFile) {
	p.Helper("sov", func(p *generator.GeneratedFile) {
		p.P(`
		func sov(x uint64) (n int) {
			return (`, p.Ident("math/bits", "Len64"), `(x | 1) + 6)/ 7
		}`)
	})
}

func HelperSOZ(p *generator.GeneratedFile) {
	p.Helper("soz", func(p *generator.GeneratedFile) {
		p.P(`func soz(x uint64) (n int) {
			return sov(uint64((x << 1) ^ uint64((int64(x) >> 63))))
		}`)
	})
}

func HelperSizeGoogleProtobufTimestamp(p *generator.GeneratedFile) {
	p.Helper("sizeGoogleProtobufTimestamp", func(p *generator.GeneratedFile) {
		p.P(`func sizeGoogleProtobufTimestamp(v *`, p.Ident("google.golang.org/protobuf/types/known/timestamppb", "Timestamp"), `) (n int) {`)
		p.P(`if v == nil {`)
		p.P(`return 0`)
		p.P(`}`)
		p.P(`if v.Seconds != 0 {`)
		p.P(`n += 1 + sov(uint64(v.Seconds))`)
		p.P(`}`)
		p.P(`if v.Nanos != 0 {`)
		p.P(`n += 1 + sov(uint64(v.Nanos))`)
		p.P(`}`)
		p.P(`return n`)
		p.P(`}`)
	})
}
