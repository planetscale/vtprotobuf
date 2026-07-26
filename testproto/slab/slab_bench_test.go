package slab

import (
	"fmt"
	"testing"
)

func BenchmarkUnmarshal(b *testing.B) {
	for _, items := range []int{1, 8, 32, 128} {
		data, err := makeBatch(items).MarshalVT()
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("stock/items=%d", items), func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := &SlabBatch{}
				if err := m.UnmarshalVT(data); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(fmt.Sprintf("slab/items=%d", items), func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := &SlabBatch{}
				if err := m.UnmarshalVTSlab(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
