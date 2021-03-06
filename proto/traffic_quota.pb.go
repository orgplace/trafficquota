// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: traffic_quota.proto

package proto

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type TakeRequest struct {
	ChunkKey   string   `protobuf:"bytes,1,opt,name=chunk_key,json=chunkKey,proto3" json:"chunk_key,omitempty"`
	BucketKeys []string `protobuf:"bytes,2,rep,name=bucket_keys,json=bucketKeys" json:"bucket_keys,omitempty"`
}

func (m *TakeRequest) Reset()         { *m = TakeRequest{} }
func (m *TakeRequest) String() string { return proto.CompactTextString(m) }
func (*TakeRequest) ProtoMessage()    {}
func (*TakeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_traffic_quota_b20b48481d300ce8, []int{0}
}
func (m *TakeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TakeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TakeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *TakeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TakeRequest.Merge(dst, src)
}
func (m *TakeRequest) XXX_Size() int {
	return m.Size()
}
func (m *TakeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TakeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TakeRequest proto.InternalMessageInfo

func (m *TakeRequest) GetChunkKey() string {
	if m != nil {
		return m.ChunkKey
	}
	return ""
}

func (m *TakeRequest) GetBucketKeys() []string {
	if m != nil {
		return m.BucketKeys
	}
	return nil
}

type TakeResponse struct {
	Allowed bool `protobuf:"varint,1,opt,name=allowed,proto3" json:"allowed,omitempty"`
}

func (m *TakeResponse) Reset()         { *m = TakeResponse{} }
func (m *TakeResponse) String() string { return proto.CompactTextString(m) }
func (*TakeResponse) ProtoMessage()    {}
func (*TakeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_traffic_quota_b20b48481d300ce8, []int{1}
}
func (m *TakeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TakeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TakeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *TakeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TakeResponse.Merge(dst, src)
}
func (m *TakeResponse) XXX_Size() int {
	return m.Size()
}
func (m *TakeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TakeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TakeResponse proto.InternalMessageInfo

func (m *TakeResponse) GetAllowed() bool {
	if m != nil {
		return m.Allowed
	}
	return false
}

func init() {
	proto.RegisterType((*TakeRequest)(nil), "orgplace.trafficquota.TakeRequest")
	proto.RegisterType((*TakeResponse)(nil), "orgplace.trafficquota.TakeResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TrafficQuotaClient is the client API for TrafficQuota service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TrafficQuotaClient interface {
	Take(ctx context.Context, in *TakeRequest, opts ...grpc.CallOption) (*TakeResponse, error)
}

type trafficQuotaClient struct {
	cc *grpc.ClientConn
}

func NewTrafficQuotaClient(cc *grpc.ClientConn) TrafficQuotaClient {
	return &trafficQuotaClient{cc}
}

func (c *trafficQuotaClient) Take(ctx context.Context, in *TakeRequest, opts ...grpc.CallOption) (*TakeResponse, error) {
	out := new(TakeResponse)
	err := c.cc.Invoke(ctx, "/orgplace.trafficquota.TrafficQuota/Take", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TrafficQuotaServer is the server API for TrafficQuota service.
type TrafficQuotaServer interface {
	Take(context.Context, *TakeRequest) (*TakeResponse, error)
}

func RegisterTrafficQuotaServer(s *grpc.Server, srv TrafficQuotaServer) {
	s.RegisterService(&_TrafficQuota_serviceDesc, srv)
}

func _TrafficQuota_Take_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TrafficQuotaServer).Take(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/orgplace.trafficquota.TrafficQuota/Take",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TrafficQuotaServer).Take(ctx, req.(*TakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TrafficQuota_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orgplace.trafficquota.TrafficQuota",
	HandlerType: (*TrafficQuotaServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Take",
			Handler:    _TrafficQuota_Take_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "traffic_quota.proto",
}

func (m *TakeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TakeRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ChunkKey) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintTrafficQuota(dAtA, i, uint64(len(m.ChunkKey)))
		i += copy(dAtA[i:], m.ChunkKey)
	}
	if len(m.BucketKeys) > 0 {
		for _, s := range m.BucketKeys {
			dAtA[i] = 0x12
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func (m *TakeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TakeResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Allowed {
		dAtA[i] = 0x8
		i++
		if m.Allowed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func encodeVarintTrafficQuota(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *TakeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ChunkKey)
	if l > 0 {
		n += 1 + l + sovTrafficQuota(uint64(l))
	}
	if len(m.BucketKeys) > 0 {
		for _, s := range m.BucketKeys {
			l = len(s)
			n += 1 + l + sovTrafficQuota(uint64(l))
		}
	}
	return n
}

func (m *TakeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Allowed {
		n += 2
	}
	return n
}

func sovTrafficQuota(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozTrafficQuota(x uint64) (n int) {
	return sovTrafficQuota(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TakeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrafficQuota
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TakeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TakeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChunkKey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrafficQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrafficQuota
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChunkKey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BucketKeys", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrafficQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrafficQuota
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BucketKeys = append(m.BucketKeys, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrafficQuota(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTrafficQuota
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TakeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrafficQuota
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TakeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TakeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Allowed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrafficQuota
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Allowed = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipTrafficQuota(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTrafficQuota
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTrafficQuota(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTrafficQuota
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
					return 0, ErrIntOverflowTrafficQuota
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTrafficQuota
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthTrafficQuota
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTrafficQuota
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipTrafficQuota(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthTrafficQuota = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTrafficQuota   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("traffic_quota.proto", fileDescriptor_traffic_quota_b20b48481d300ce8) }

var fileDescriptor_traffic_quota_b20b48481d300ce8 = []byte{
	// 227 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2e, 0x29, 0x4a, 0x4c,
	0x4b, 0xcb, 0x4c, 0x8e, 0x2f, 0x2c, 0xcd, 0x2f, 0x49, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x12, 0xcd, 0x2f, 0x4a, 0x2f, 0xc8, 0x49, 0x4c, 0x4e, 0xd5, 0x83, 0xca, 0x82, 0x25, 0x95, 0xbc,
	0xb9, 0xb8, 0x43, 0x12, 0xb3, 0x53, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xa4, 0xb9,
	0x38, 0x93, 0x33, 0x4a, 0xf3, 0xb2, 0xe3, 0xb3, 0x53, 0x2b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38,
	0x83, 0x38, 0xc0, 0x02, 0xde, 0xa9, 0x95, 0x42, 0xf2, 0x5c, 0xdc, 0x49, 0xa5, 0xc9, 0xd9, 0xa9,
	0x25, 0x20, 0xd9, 0x62, 0x09, 0x26, 0x05, 0x66, 0x0d, 0xce, 0x20, 0x2e, 0x88, 0x90, 0x77, 0x6a,
	0x65, 0xb1, 0x92, 0x06, 0x17, 0x0f, 0xc4, 0xb0, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x09,
	0x2e, 0xf6, 0xc4, 0x9c, 0x9c, 0xfc, 0xf2, 0xd4, 0x14, 0xb0, 0x59, 0x1c, 0x41, 0x30, 0xae, 0x51,
	0x22, 0x17, 0x4f, 0x08, 0xc4, 0x19, 0x81, 0x20, 0x67, 0x08, 0x05, 0x72, 0xb1, 0x80, 0x74, 0x0a,
	0x29, 0xe9, 0x61, 0x75, 0xa6, 0x1e, 0x92, 0x1b, 0xa5, 0x94, 0xf1, 0xaa, 0x81, 0x58, 0xad, 0xc4,
	0xe0, 0x24, 0x7f, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e,
	0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0xac, 0xe0, 0x10,
	0x49, 0x62, 0x03, 0x53, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4a, 0xb9, 0x70, 0x63, 0x2f,
	0x01, 0x00, 0x00,
}
