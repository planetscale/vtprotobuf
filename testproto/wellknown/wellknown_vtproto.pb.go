// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// protoc-gen-go-vtproto version: (devel)
// source: wellknown/wellknown.proto

package wellknown

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	bits "math/bits"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

func (m *WellKnownFields) CloneVT() *WellKnownFields {
	if m == nil {
		return (*WellKnownFields)(nil)
	}
	r := &WellKnownFields{}
	if rhs := m.Timestamp; rhs != nil {
		if vtpb, ok := interface{}(rhs).(interface{ CloneVT() *timestamppb.Timestamp }); ok {
			r.Timestamp = vtpb.CloneVT()
		} else {
			r.Timestamp = proto.Clone(rhs).(*timestamppb.Timestamp)
		}
	}
	if rhs := m.Duration; rhs != nil {
		if vtpb, ok := interface{}(rhs).(interface{ CloneVT() *durationpb.Duration }); ok {
			r.Duration = vtpb.CloneVT()
		} else {
			r.Duration = proto.Clone(rhs).(*durationpb.Duration)
		}
	}
	if m.OneofField != nil {
		r.OneofField = m.OneofField.(interface {
			CloneVT() isWellKnownFields_OneofField
		}).CloneVT()
	}
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *WellKnownFields) CloneMessageVT() proto.Message {
	return m.CloneVT()
}

func (m *WellKnownFields_OneofTimestamp) CloneVT() isWellKnownFields_OneofField {
	if m == nil {
		return (*WellKnownFields_OneofTimestamp)(nil)
	}
	r := &WellKnownFields_OneofTimestamp{}
	if rhs := m.OneofTimestamp; rhs != nil {
		if vtpb, ok := interface{}(rhs).(interface{ CloneVT() *timestamppb.Timestamp }); ok {
			r.OneofTimestamp = vtpb.CloneVT()
		} else {
			r.OneofTimestamp = proto.Clone(rhs).(*timestamppb.Timestamp)
		}
	}
	return r
}

func (m *WellKnownFields_OneofDuration) CloneVT() isWellKnownFields_OneofField {
	if m == nil {
		return (*WellKnownFields_OneofDuration)(nil)
	}
	r := &WellKnownFields_OneofDuration{}
	if rhs := m.OneofDuration; rhs != nil {
		if vtpb, ok := interface{}(rhs).(interface{ CloneVT() *durationpb.Duration }); ok {
			r.OneofDuration = vtpb.CloneVT()
		} else {
			r.OneofDuration = proto.Clone(rhs).(*durationpb.Duration)
		}
	}
	return r
}

func (this *WellKnownFields) EqualVT(that *WellKnownFields) bool {
	if this == that {
		return true
	} else if this == nil || that == nil {
		return false
	}
	if this.OneofField == nil && that.OneofField != nil {
		return false
	} else if this.OneofField != nil {
		if that.OneofField == nil {
			return false
		}
		if !this.OneofField.(interface {
			EqualVT(isWellKnownFields_OneofField) bool
		}).EqualVT(that.OneofField) {
			return false
		}
	}
	if equal, ok := interface{}(this.Timestamp).(interface {
		EqualVT(*timestamppb.Timestamp) bool
	}); ok {
		if !equal.EqualVT(that.Timestamp) {
			return false
		}
	} else if !proto.Equal(this.Timestamp, that.Timestamp) {
		return false
	}
	if equal, ok := interface{}(this.Duration).(interface {
		EqualVT(*durationpb.Duration) bool
	}); ok {
		if !equal.EqualVT(that.Duration) {
			return false
		}
	} else if !proto.Equal(this.Duration, that.Duration) {
		return false
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *WellKnownFields) EqualMessageVT(thatMsg proto.Message) bool {
	that, ok := thatMsg.(*WellKnownFields)
	if !ok {
		return false
	}
	return this.EqualVT(that)
}
func (this *WellKnownFields_OneofTimestamp) EqualVT(thatIface isWellKnownFields_OneofField) bool {
	that, ok := thatIface.(*WellKnownFields_OneofTimestamp)
	if !ok {
		return false
	}
	if this == that {
		return true
	}
	if this == nil && that != nil || this != nil && that == nil {
		return false
	}
	if p, q := this.OneofTimestamp, that.OneofTimestamp; p != q {
		if p == nil {
			p = &timestamppb.Timestamp{}
		}
		if q == nil {
			q = &timestamppb.Timestamp{}
		}
		if equal, ok := interface{}(p).(interface {
			EqualVT(*timestamppb.Timestamp) bool
		}); ok {
			if !equal.EqualVT(q) {
				return false
			}
		} else if !proto.Equal(p, q) {
			return false
		}
	}
	return true
}

func (this *WellKnownFields_OneofDuration) EqualVT(thatIface isWellKnownFields_OneofField) bool {
	that, ok := thatIface.(*WellKnownFields_OneofDuration)
	if !ok {
		return false
	}
	if this == that {
		return true
	}
	if this == nil && that != nil || this != nil && that == nil {
		return false
	}
	if p, q := this.OneofDuration, that.OneofDuration; p != q {
		if p == nil {
			p = &durationpb.Duration{}
		}
		if q == nil {
			q = &durationpb.Duration{}
		}
		if equal, ok := interface{}(p).(interface {
			EqualVT(*durationpb.Duration) bool
		}); ok {
			if !equal.EqualVT(q) {
				return false
			}
		} else if !proto.Equal(p, q) {
			return false
		}
	}
	return true
}

func (m *WellKnownFields) MarshalVT() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVT(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WellKnownFields) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *WellKnownFields) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if vtmsg, ok := m.OneofField.(interface {
		MarshalToSizedBufferVT([]byte) (int, error)
	}); ok {
		size, err := vtmsg.MarshalToSizedBufferVT(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	if m.Duration != nil {
		size, err := marshal_durationpb_Duration(dAtA[:i], m.Duration)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x12
	}
	if m.Timestamp != nil {
		size, err := marshal_timestamppb_Timestamp(dAtA[:i], m.Timestamp)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *WellKnownFields_OneofTimestamp) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *WellKnownFields_OneofTimestamp) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OneofTimestamp != nil {
		size, err := marshal_timestamppb_Timestamp(dAtA[:i], m.OneofTimestamp)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x1a
	}
	return len(dAtA) - i, nil
}
func (m *WellKnownFields_OneofDuration) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *WellKnownFields_OneofDuration) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OneofDuration != nil {
		size, err := marshal_durationpb_Duration(dAtA[:i], m.OneofDuration)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x22
	}
	return len(dAtA) - i, nil
}
func marshal_durationpb_Duration(dAtA []byte, m *durationpb.Duration) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Nanos != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Nanos))
		i--
		dAtA[i] = 0x10
	}
	if m.Seconds != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Seconds))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func marshal_timestamppb_Timestamp(dAtA []byte, m *timestamppb.Timestamp) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Nanos != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Nanos))
		i--
		dAtA[i] = 0x10
	}
	if m.Seconds != 0 {
		i = encodeVarint(dAtA, i, uint64(m.Seconds))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func sov(x uint64) (n int) {
	return (bits.Len64(x|1) + 6) / 7
}
func encodeVarint(dAtA []byte, offset int, v uint64) int {
	offset -= sov(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *WellKnownFields) MarshalVTStrict() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVTStrict(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WellKnownFields) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *WellKnownFields) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if msg, ok := m.OneofField.(*WellKnownFields_OneofDuration); ok {
		size, err := msg.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	if msg, ok := m.OneofField.(*WellKnownFields_OneofTimestamp); ok {
		size, err := msg.MarshalToSizedBufferVTStrict(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
	}
	if m.Duration != nil {
		size, err := marshal_durationpb_Duration(dAtA[:i], m.Duration)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x12
	}
	if m.Timestamp != nil {
		size, err := marshal_timestamppb_Timestamp(dAtA[:i], m.Timestamp)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *WellKnownFields_OneofTimestamp) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *WellKnownFields_OneofTimestamp) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OneofTimestamp != nil {
		size, err := marshal_timestamppb_Timestamp(dAtA[:i], m.OneofTimestamp)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x1a
	}
	return len(dAtA) - i, nil
}
func (m *WellKnownFields_OneofDuration) MarshalToVTStrict(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVTStrict(dAtA[:size])
}

func (m *WellKnownFields_OneofDuration) MarshalToSizedBufferVTStrict(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.OneofDuration != nil {
		size, err := marshal_durationpb_Duration(dAtA[:i], m.OneofDuration)
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarint(dAtA, i, uint64(size))
		i--
		dAtA[i] = 0x22
	}
	return len(dAtA) - i, nil
}
func (m *WellKnownFields) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Timestamp != nil {
		l = size_timestamppb_Timestamp(m.Timestamp)
		n += 1 + l + sov(uint64(l))
	}
	if m.Duration != nil {
		l = size_durationpb_Duration(m.Duration)
		n += 1 + l + sov(uint64(l))
	}
	if vtmsg, ok := m.OneofField.(interface{ SizeVT() int }); ok {
		n += vtmsg.SizeVT()
	}
	n += len(m.unknownFields)
	return n
}

func (m *WellKnownFields_OneofTimestamp) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.OneofTimestamp != nil {
		l = size_timestamppb_Timestamp(m.OneofTimestamp)
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func (m *WellKnownFields_OneofDuration) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.OneofDuration != nil {
		l = size_durationpb_Duration(m.OneofDuration)
		n += 1 + l + sov(uint64(l))
	}
	return n
}
func size_timestamppb_Timestamp(m *timestamppb.Timestamp) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Seconds != 0 {
		n += 1 + sov(uint64(m.Seconds))
	}
	if m.Nanos != 0 {
		n += 1 + sov(uint64(m.Nanos))
	}
	return n
}

func size_durationpb_Duration(m *durationpb.Duration) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Seconds != 0 {
		n += 1 + sov(uint64(m.Seconds))
	}
	if m.Nanos != 0 {
		n += 1 + sov(uint64(m.Nanos))
	}
	return n
}

func soz(x uint64) (n int) {
	return sov(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *WellKnownFields) UnmarshalVT(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WellKnownFields: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WellKnownFields: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &timestamppb.Timestamp{}
			}
			if unmarshal, ok := interface{}(m.Timestamp).(interface {
				UnmarshalVT([]byte) error
			}); ok {
				if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
			} else {
				if err := proto.Unmarshal(dAtA[iNdEx:postIndex], m.Timestamp); err != nil {
					return err
				}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Duration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Duration == nil {
				m.Duration = &durationpb.Duration{}
			}
			if unmarshal, ok := interface{}(m.Duration).(interface {
				UnmarshalVT([]byte) error
			}); ok {
				if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
					return err
				}
			} else {
				if err := proto.Unmarshal(dAtA[iNdEx:postIndex], m.Duration); err != nil {
					return err
				}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OneofTimestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if oneof, ok := m.OneofField.(*WellKnownFields_OneofTimestamp); ok {
				if unmarshal, ok := interface{}(oneof.OneofTimestamp).(interface {
					UnmarshalVT([]byte) error
				}); ok {
					if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
						return err
					}
				} else {
					if err := proto.Unmarshal(dAtA[iNdEx:postIndex], oneof.OneofTimestamp); err != nil {
						return err
					}
				}
			} else {
				v := &timestamppb.Timestamp{}
				if unmarshal, ok := interface{}(v).(interface {
					UnmarshalVT([]byte) error
				}); ok {
					if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
						return err
					}
				} else {
					if err := proto.Unmarshal(dAtA[iNdEx:postIndex], v); err != nil {
						return err
					}
				}
				m.OneofField = &WellKnownFields_OneofTimestamp{OneofTimestamp: v}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OneofDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLength
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLength
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if oneof, ok := m.OneofField.(*WellKnownFields_OneofDuration); ok {
				if unmarshal, ok := interface{}(oneof.OneofDuration).(interface {
					UnmarshalVT([]byte) error
				}); ok {
					if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
						return err
					}
				} else {
					if err := proto.Unmarshal(dAtA[iNdEx:postIndex], oneof.OneofDuration); err != nil {
						return err
					}
				}
			} else {
				v := &durationpb.Duration{}
				if unmarshal, ok := interface{}(v).(interface {
					UnmarshalVT([]byte) error
				}); ok {
					if err := unmarshal.UnmarshalVT(dAtA[iNdEx:postIndex]); err != nil {
						return err
					}
				} else {
					if err := proto.Unmarshal(dAtA[iNdEx:postIndex], v); err != nil {
						return err
					}
				}
				m.OneofField = &WellKnownFields_OneofDuration{OneofDuration: v}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.unknownFields = append(m.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func skip(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflow
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflow
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLength
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroup
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLength
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLength        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflow          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroup = fmt.Errorf("proto: unexpected end of group")
)
