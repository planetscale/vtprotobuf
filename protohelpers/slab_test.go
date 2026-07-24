package protohelpers

import (
	"testing"

	"google.golang.org/protobuf/encoding/protowire"
)

func TestSlabNextHandsOutDistinctZeroedElements(t *testing.T) {
	var s Slab[int64]
	seen := make(map[*int64]bool)
	for i := 0; i < 3000; i++ {
		p := s.Next()
		if *p != 0 {
			t.Fatalf("element %d not zeroed: %d", i, *p)
		}
		if seen[p] {
			t.Fatalf("element %d handed out twice", i)
		}
		seen[p] = true
		*p = int64(i)
	}
	if len(seen) != 3000 {
		t.Fatalf("expected 3000 distinct elements, got %d", len(seen))
	}
}

func TestSlabNextValue(t *testing.T) {
	var s Slab[string]
	p := s.NextValue("hello")
	if *p != "hello" {
		t.Fatalf("NextValue stored %q", *p)
	}
	q := s.NextValue("world")
	if *p != "hello" || *q != "world" {
		t.Fatalf("NextValue overwrote a previous element: %q %q", *p, *q)
	}
}

func TestSlabReserveExact(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		var s Slab[uint64]
		s.Reserve(100)
		for i := 0; i < 100; i++ {
			*s.Next() = uint64(i)
		}
	})
	if allocs != 1 {
		t.Fatalf("expected exactly 1 chunk allocation for a reserved slab, got %v", allocs)
	}
}

func TestSlabReserveUndershootFallsBackToGrowth(t *testing.T) {
	var s Slab[int32]
	s.Reserve(4)
	seen := make(map[*int32]bool)
	for i := 0; i < 100; i++ {
		p := s.Next()
		if seen[p] {
			t.Fatal("element handed out twice after reservation undershoot")
		}
		seen[p] = true
	}
}

func TestSlabReserveAfterUseIsNoop(t *testing.T) {
	var s Slab[int32]
	p := s.Next()
	*p = 42
	s.Reserve(1000)
	if *p != 42 {
		t.Fatal("Reserve after use disturbed an existing element")
	}
}

func TestCountFields(t *testing.T) {
	var data []byte
	data = protowire.AppendTag(data, 1, protowire.VarintType)
	data = protowire.AppendVarint(data, 300)
	data = protowire.AppendTag(data, 2, protowire.BytesType)
	data = protowire.AppendBytes(data, []byte("abc"))
	data = protowire.AppendTag(data, 3, protowire.Fixed64Type)
	data = protowire.AppendFixed64(data, 7)
	data = protowire.AppendTag(data, 2, protowire.BytesType)
	data = protowire.AppendBytes(data, []byte("def"))
	data = protowire.AppendTag(data, 4, protowire.Fixed32Type)
	data = protowire.AppendFixed32(data, 7)
	data = protowire.AppendTag(data, 5, protowire.BytesType)
	// a nested message containing field 2 must NOT be counted: the scan is
	// top-level only
	nested := protowire.AppendTag(nil, 2, protowire.BytesType)
	nested = protowire.AppendBytes(nested, []byte("inner"))
	data = protowire.AppendBytes(data, nested)
	data = protowire.AppendTag(data, 5, protowire.BytesType)
	data = protowire.AppendBytes(data, nil)

	fieldNums := []int32{2, 5, 9}
	counts := make([]int, 3)
	CountFields(data, fieldNums, counts)
	if counts[0] != 2 || counts[1] != 2 || counts[2] != 0 {
		t.Fatalf("got counts %v, want [2 2 0]", counts)
	}

	// varint field 2 must not be counted (only length-delimited occurrences)
	data2 := protowire.AppendTag(nil, 2, protowire.VarintType)
	data2 = protowire.AppendVarint(data2, 1)
	counts2 := make([]int, 1)
	CountFields(data2, []int32{2}, counts2)
	if counts2[0] != 0 {
		t.Fatalf("varint occurrence counted: %v", counts2)
	}
}

func TestCountFieldsStopsAtGroups(t *testing.T) {
	var data []byte
	data = protowire.AppendTag(data, 2, protowire.BytesType)
	data = protowire.AppendBytes(data, []byte("abc"))
	data = protowire.AppendTag(data, 3, protowire.StartGroupType)
	data = protowire.AppendTag(data, 2, protowire.BytesType) // inside the group
	data = protowire.AppendBytes(data, []byte("x"))
	data = protowire.AppendTag(data, 3, protowire.EndGroupType)
	data = protowire.AppendTag(data, 2, protowire.BytesType) // after the group
	data = protowire.AppendBytes(data, []byte("def"))

	counts := make([]int, 1)
	CountFields(data, []int32{2}, counts)
	// the scan bails at the group: only the occurrence before it is counted
	if counts[0] != 1 {
		t.Fatalf("got count %d, want 1 (scan should stop at the group)", counts[0])
	}
}

func TestCountFieldsMalformedInput(t *testing.T) {
	var data []byte
	data = protowire.AppendTag(data, 2, protowire.BytesType)
	data = protowire.AppendBytes(data, []byte("abcdef"))
	data = protowire.AppendTag(data, 2, protowire.BytesType)
	data = protowire.AppendBytes(data, []byte("ghijkl"))

	// truncations at every offset must never panic and never overcount
	for i := 0; i <= len(data); i++ {
		counts := make([]int, 1)
		CountFields(data[:i], []int32{2}, counts)
		if counts[0] > 2 {
			t.Fatalf("truncation at %d overcounted: %d", i, counts[0])
		}
	}

	// a length prefix pointing far past the buffer must not panic
	bad := protowire.AppendTag(nil, 2, protowire.BytesType)
	bad = protowire.AppendVarint(bad, 1<<40)
	counts := make([]int, 1)
	CountFields(bad, []int32{2}, counts)
}
