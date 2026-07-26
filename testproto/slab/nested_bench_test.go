package slab

import (
	"fmt"
	"testing"
)

// BenchmarkUnmarshalNested compares stock and slab decoding on a deeply
// nested document-style payload (four to five message levels, with maps,
// oneofs and optional scalars at every level).
func BenchmarkUnmarshalNested(b *testing.B) {
	for _, lines := range []int{1, 8, 32, 128} {
		data, err := makeNestedOrder(lines).MarshalVT()
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("stock/lines=%d", lines), func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := &NestedOrder{}
				if err := m.UnmarshalVT(data); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(fmt.Sprintf("slab/lines=%d", lines), func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := &NestedOrder{}
				if err := m.UnmarshalVTSlab(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
