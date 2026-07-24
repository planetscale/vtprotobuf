package protohelpers

// SlabUnmarshalThreshold is the payload size (in bytes) below which the
// generated UnmarshalVTSlab methods fall back to plain UnmarshalVT: for tiny
// messages the arena setup and counting pre-pass cost more than the handful
// of allocations they would save.
const SlabUnmarshalThreshold = 256

// Slab hands out zeroed *T elements carved from chunked backing arrays, so N
// element allocations collapse into O(log N) chunk allocations. Elements are
// handed out exactly once and never reused, so pointers into a chunk stay
// valid forever; a retained pointer pins its whole chunk, which is the
// intended trade-off for RPC-lifetime messages.
type Slab[T any] struct {
	chunk []T
	// grow is the size of the next growth chunk. It intentionally does not
	// account for Reserve: when a reservation is undershot, growth restarts
	// small instead of doubling the (possibly large) reserved chunk, so an
	// overflow of a few elements never allocates reservation-sized chunks.
	grow int
}

// Next returns a pointer to a new zeroed element.
func (s *Slab[T]) Next() *T {
	if len(s.chunk) == cap(s.chunk) {
		n := s.grow
		switch {
		case n == 0:
			n = 8
		case n > 1024:
			n = 1024
		}
		s.chunk = make([]T, 0, n)
		s.grow = n * 2
	}
	s.chunk = s.chunk[:len(s.chunk)+1]
	return &s.chunk[len(s.chunk)-1]
}

// NextValue returns a pointer to a new element initialized to v.
func (s *Slab[T]) NextValue(v T) *T {
	p := s.Next()
	*p = v
	return p
}

// Reserve pre-sizes an empty slab so that the next n calls to Next land in
// one exactly-sized chunk. It has no effect once the slab has a chunk, and a
// reservation that turns out too small simply falls back to chunked growth,
// so callers may pass a best-effort estimate.
func (s *Slab[T]) Reserve(n int) {
	if s.chunk == nil && n > 0 {
		s.chunk = make([]T, 0, n)
	}
}

// CountFields performs one linear pass over the top-level records of a
// serialized message and increments counts[i] for every length-delimited
// occurrence of field number fieldNums[i]. fieldNums and counts must have the
// same length. The counts are used to Reserve exactly-sized slabs before
// decoding; on malformed input (or on encountering a group, which would
// require nested tracking) it returns early with the counts gathered so far,
// which at worst under-reserves and falls back to chunked growth.
func CountFields(data []byte, fieldNums []int32, counts []int) {
	l := len(data)
	i := 0
	for i < l {
		key, n := countVarint(data, i)
		if n < 0 {
			return
		}
		i = n
		switch key & 0x7 {
		case 0: // varint
			_, n = countVarint(data, i)
			if n < 0 {
				return
			}
			i = n
		case 1: // fixed64
			i += 8
		case 2: // length-delimited
			length, n := countVarint(data, i)
			if n < 0 {
				return
			}
			i = n
			fieldNum := int32(key >> 3)
			for j, want := range fieldNums {
				if want == fieldNum {
					counts[j]++
					break
				}
			}
			i += int(length)
			if i < 0 {
				return
			}
		case 5: // fixed32
			i += 4
		default: // groups or invalid wire type
			return
		}
	}
}

// countVarint decodes a varint at data[i:] and returns its value and the
// offset of the next byte, or a negative offset on truncated input.
func countVarint(data []byte, i int) (v uint64, next int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if i >= len(data) {
			return 0, -1
		}
		b := data[i]
		i++
		v |= uint64(b&0x7F) << shift
		if b < 0x80 {
			return v, i
		}
	}
	return 0, -1
}
