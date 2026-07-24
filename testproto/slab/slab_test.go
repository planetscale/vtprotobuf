package slab

import (
	"fmt"
	"testing"

	grpccodec "github.com/planetscale/vtprotobuf/codec/grpc"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
)

func ptr[T any](v T) *T { return &v }

func makeItem(i int, withChildren bool) *SlabItem {
	item := &SlabItem{
		Key:   fmt.Sprintf("/keys/%016d", i),
		Value: []byte(fmt.Sprintf("value payload %d padding padding padding", i)),
		Meta:  &SlabMeta{CreatedAt: int64(1000 + i)},
	}
	// exercise both set and unset presence fields
	if i%2 == 0 {
		item.ExpectedVersion = ptr(int64(i))
		item.ClientIdentity = ptr(fmt.Sprintf("client-%d", i))
		item.Flag = ptr(i%4 == 0)
		item.Score = ptr(float64(i) * 1.5)
	}
	if i%3 == 0 {
		item.Stamp = ptr(uint64(i) << 20)
		item.Flavor = ptr(SlabFlavor(i % 3))
		item.Meta.Ttl = ptr(uint32(i))
		item.Meta.Weight = ptr(float32(i) / 3)
	}
	if withChildren && i%5 == 0 {
		item.Children = []*SlabItem{makeItem(10*i+1, false), makeItem(10*i+2, false)}
	}
	return item
}

func makeBatch(items int) *SlabBatch {
	batch := &SlabBatch{
		Shard:     ptr(int64(42)),
		Meta:      &SlabMeta{CreatedAt: 7, Ttl: ptr(uint32(99))},
		PackedIds: []int64{},
	}
	for i := 0; i < items; i++ {
		batch.Items = append(batch.Items, makeItem(i, true))
		batch.PackedIds = append(batch.PackedIds, int64(i)*7)
	}
	for i := 0; i < items/2; i++ {
		batch.Ranges = append(batch.Ranges, &SlabRange{
			StartInclusive: fmt.Sprintf("/start/%d", i),
			EndExclusive:   fmt.Sprintf("/end/%d", i),
			Delta:          ptr(int64(-i)),
		})
	}
	if items > 0 {
		batch.Index = map[string]*SlabItem{}
		for i := 0; i < (items+3)/4; i++ {
			batch.Index[fmt.Sprintf("idx-%d", i)] = makeItem(1000+i, false)
		}
		batch.Extra = &SlabBatch_ExtraItem{ExtraItem: makeItem(9999, false)}
	} else {
		batch.Extra = &SlabBatch_Note{Note: "empty"}
	}
	return batch
}

func TestUnmarshalVTSlabEquivalence(t *testing.T) {
	for _, items := range []int{0, 1, 3, 8, 32, 200} {
		t.Run(fmt.Sprintf("items=%d", items), func(t *testing.T) {
			orig := makeBatch(items)
			data, err := orig.MarshalVT()
			if err != nil {
				t.Fatal(err)
			}

			stock := &SlabBatch{}
			if err := stock.UnmarshalVT(data); err != nil {
				t.Fatal(err)
			}
			slab := &SlabBatch{}
			if err := slab.UnmarshalVTSlab(data); err != nil {
				t.Fatal(err)
			}

			if !proto.Equal(orig, slab) {
				t.Fatalf("slab-decoded message differs from original:\norig: %v\nslab: %v", orig, slab)
			}
			if !proto.Equal(stock, slab) {
				t.Fatalf("slab-decoded message differs from stock-decoded one")
			}
		})
	}
}

// Small payloads take the bypass branch and must behave identically.
func TestUnmarshalVTSlabSmallPayloadBypass(t *testing.T) {
	orig := &SlabBatch{Shard: ptr(int64(1)), Items: []*SlabItem{{Key: "k", Value: []byte("v")}}}
	data, err := orig.MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) >= 256 {
		t.Fatalf("test payload unexpectedly large: %d bytes", len(data))
	}
	slab := &SlabBatch{}
	if err := slab.UnmarshalVTSlab(data); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(orig, slab) {
		t.Fatal("bypass path decoded a different message")
	}
}

func TestUnmarshalVTSlabRetainsUnknownFields(t *testing.T) {
	data, err := makeBatch(8).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	// append unknown fields of several wire types
	unknown := protowire.AppendTag(nil, 999, protowire.BytesType)
	unknown = protowire.AppendBytes(unknown, []byte("unknown payload"))
	unknown = protowire.AppendTag(unknown, 1000, protowire.VarintType)
	unknown = protowire.AppendVarint(unknown, 12345)
	data = append(data, unknown...)

	slab := &SlabBatch{}
	if err := slab.UnmarshalVTSlab(data); err != nil {
		t.Fatal(err)
	}
	stock := &SlabBatch{}
	if err := stock.UnmarshalVT(data); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(stock, slab) {
		t.Fatal("slab decoding with unknown fields differs from stock")
	}

	reData, err := slab.MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	reDecoded := &SlabBatch{}
	if err := reDecoded.UnmarshalVT(reData); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(stock, reDecoded) {
		t.Fatal("unknown fields were not retained through a slab decode + re-encode")
	}
}

// The decoded message must not alias the input buffer.
func TestUnmarshalVTSlabNoAliasing(t *testing.T) {
	orig := makeBatch(16)
	data, err := orig.MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	slab := &SlabBatch{}
	if err := slab.UnmarshalVTSlab(data); err != nil {
		t.Fatal(err)
	}
	for i := range data {
		data[i] = 0xFF
	}
	if !proto.Equal(orig, slab) {
		t.Fatal("decoded message changed when the source buffer was overwritten")
	}
}

// Unmarshalling into a non-fresh message follows proto merge semantics;
// slab decoding must match stock decoding exactly.
func TestUnmarshalVTSlabMerge(t *testing.T) {
	first, err := makeBatch(8).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	second, err := makeBatch(5).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}

	stock := &SlabBatch{}
	slab := &SlabBatch{}
	for _, data := range [][]byte{first, second} {
		if err := stock.UnmarshalVT(data); err != nil {
			t.Fatal(err)
		}
		if err := slab.UnmarshalVTSlab(data); err != nil {
			t.Fatal(err)
		}
	}
	if !proto.Equal(stock, slab) {
		t.Fatal("merge semantics differ between stock and slab decoding")
	}
}

// Truncations and corruptions must never panic (this exercises both the
// counting pre-pass and the decode loops on malformed input).
func TestUnmarshalVTSlabMalformedInput(t *testing.T) {
	data, err := makeBatch(8).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 256 {
		t.Fatalf("payload must exceed the bypass threshold, got %d bytes", len(data))
	}
	for i := 0; i <= len(data); i++ {
		m := &SlabBatch{}
		_ = m.UnmarshalVTSlab(data[:i]) // must not panic
	}
	for i := 0; i < len(data); i++ {
		corrupted := append([]byte(nil), data...)
		corrupted[i] ^= 0xA5
		m := &SlabBatch{}
		_ = m.UnmarshalVTSlab(corrupted) // must not panic
	}
}

// The whole point: slab decoding must allocate substantially less than
// stock decoding on repeated-message-heavy payloads.
func TestUnmarshalVTSlabReducesAllocations(t *testing.T) {
	data, err := makeBatch(32).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	stockAllocs := testing.AllocsPerRun(50, func() {
		m := &SlabBatch{}
		if err := m.UnmarshalVT(data); err != nil {
			t.Fatal(err)
		}
	})
	slabAllocs := testing.AllocsPerRun(50, func() {
		m := &SlabBatch{}
		if err := m.UnmarshalVTSlab(data); err != nil {
			t.Fatal(err)
		}
	})
	t.Logf("allocs/op: stock=%v slab=%v", stockAllocs, slabAllocs)
	if slabAllocs >= stockAllocs {
		t.Fatalf("slab decoding did not reduce allocations: stock=%v slab=%v", stockAllocs, slabAllocs)
	}
}

// The bundled gRPC codec must transparently dispatch to UnmarshalVTSlab for
// opted-in messages and keep using UnmarshalVT for everything else.
func TestGRPCCodecPrefersSlab(t *testing.T) {
	var codec grpccodec.Codec

	orig := makeBatch(16)
	data, err := codec.Marshal(orig)
	if err != nil {
		t.Fatal(err)
	}
	decoded := &SlabBatch{}
	if err := codec.Unmarshal(data, decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(orig, decoded) {
		t.Fatal("codec round-trip through the slab entry point lost data")
	}

	plain := &SlabPlain{Name: "n"}
	plainData, err := codec.Marshal(plain)
	if err != nil {
		t.Fatal(err)
	}
	plainDecoded := &SlabPlain{}
	if err := codec.Unmarshal(plainData, plainDecoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(plain, plainDecoded) {
		t.Fatal("codec round-trip of a non-slab message lost data")
	}
}
