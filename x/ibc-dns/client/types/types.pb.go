// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ibc/dns/client/types.proto

package types

import (
	context "context"
	fmt "fmt"
	types "github.com/datachainlab/cosmos-sdk-interchain-dns/x/ibc-dns/common/types"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MsgRegisterDomain struct {
	// the port on which the packet will be sent
	SourcePort string `protobuf:"bytes,1,opt,name=source_port,json=sourcePort,proto3" json:"source_port,omitempty" yaml:"source_port"`
	// the channel by which the packet will be sent
	SourceChannel string `protobuf:"bytes,2,opt,name=source_channel,json=sourceChannel,proto3" json:"source_channel,omitempty" yaml:"source_channel"`
	Domain        string `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Metadata      []byte `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Sender        []byte `protobuf:"bytes,5,opt,name=sender,proto3" json:"sender,omitempty"`
}

func (m *MsgRegisterDomain) Reset()         { *m = MsgRegisterDomain{} }
func (m *MsgRegisterDomain) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterDomain) ProtoMessage()    {}
func (*MsgRegisterDomain) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd45cdcc41e691b0, []int{0}
}
func (m *MsgRegisterDomain) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRegisterDomain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterDomain.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRegisterDomain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterDomain.Merge(m, src)
}
func (m *MsgRegisterDomain) XXX_Size() int {
	return m.Size()
}
func (m *MsgRegisterDomain) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRegisterDomain.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRegisterDomain proto.InternalMessageInfo

type MsgRegisterDomainResponse struct {
}

func (m *MsgRegisterDomainResponse) Reset()         { *m = MsgRegisterDomainResponse{} }
func (m *MsgRegisterDomainResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRegisterDomainResponse) ProtoMessage()    {}
func (*MsgRegisterDomainResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd45cdcc41e691b0, []int{1}
}
func (m *MsgRegisterDomainResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRegisterDomainResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRegisterDomainResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRegisterDomainResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRegisterDomainResponse.Merge(m, src)
}
func (m *MsgRegisterDomainResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgRegisterDomainResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRegisterDomainResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRegisterDomainResponse proto.InternalMessageInfo

type MsgDomainAssociationCreate struct {
	Sender    []byte             `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	DnsId     types.LocalDNSID   `protobuf:"bytes,2,opt,name=dns_id,json=dnsId,proto3" json:"dns_id"`
	SrcClient types.ClientDomain `protobuf:"bytes,3,opt,name=src_client,json=srcClient,proto3" json:"src_client"`
	DstClient types.ClientDomain `protobuf:"bytes,4,opt,name=dst_client,json=dstClient,proto3" json:"dst_client"`
}

func (m *MsgDomainAssociationCreate) Reset()         { *m = MsgDomainAssociationCreate{} }
func (m *MsgDomainAssociationCreate) String() string { return proto.CompactTextString(m) }
func (*MsgDomainAssociationCreate) ProtoMessage()    {}
func (*MsgDomainAssociationCreate) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd45cdcc41e691b0, []int{2}
}
func (m *MsgDomainAssociationCreate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDomainAssociationCreate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDomainAssociationCreate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDomainAssociationCreate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDomainAssociationCreate.Merge(m, src)
}
func (m *MsgDomainAssociationCreate) XXX_Size() int {
	return m.Size()
}
func (m *MsgDomainAssociationCreate) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDomainAssociationCreate.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDomainAssociationCreate proto.InternalMessageInfo

type MsgDomainAssociationCreateResponse struct {
}

func (m *MsgDomainAssociationCreateResponse) Reset()         { *m = MsgDomainAssociationCreateResponse{} }
func (m *MsgDomainAssociationCreateResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDomainAssociationCreateResponse) ProtoMessage()    {}
func (*MsgDomainAssociationCreateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd45cdcc41e691b0, []int{3}
}
func (m *MsgDomainAssociationCreateResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDomainAssociationCreateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDomainAssociationCreateResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDomainAssociationCreateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDomainAssociationCreateResponse.Merge(m, src)
}
func (m *MsgDomainAssociationCreateResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgDomainAssociationCreateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDomainAssociationCreateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDomainAssociationCreateResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgRegisterDomain)(nil), "ibc.dns.client.MsgRegisterDomain")
	proto.RegisterType((*MsgRegisterDomainResponse)(nil), "ibc.dns.client.MsgRegisterDomainResponse")
	proto.RegisterType((*MsgDomainAssociationCreate)(nil), "ibc.dns.client.MsgDomainAssociationCreate")
	proto.RegisterType((*MsgDomainAssociationCreateResponse)(nil), "ibc.dns.client.MsgDomainAssociationCreateResponse")
}

func init() { proto.RegisterFile("ibc/dns/client/types.proto", fileDescriptor_cd45cdcc41e691b0) }

var fileDescriptor_cd45cdcc41e691b0 = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x69, 0x1a, 0xd1, 0x0d, 0x44, 0x62, 0x05, 0x25, 0x35, 0xc8, 0x29, 0x16, 0x87, 0x82,
	0x14, 0x5b, 0x0a, 0x87, 0x4a, 0x3d, 0xd1, 0x24, 0x07, 0x2a, 0x11, 0x84, 0xcc, 0x8d, 0x03, 0xd1,
	0x7a, 0x77, 0xe5, 0xac, 0x88, 0x77, 0xa3, 0x9d, 0xad, 0x44, 0xfe, 0x80, 0x23, 0x9f, 0xc0, 0xe7,
	0xf4, 0xd8, 0x23, 0xa7, 0x28, 0x4a, 0xfe, 0xa0, 0x37, 0x6e, 0xc8, 0xbb, 0x4e, 0x49, 0x53, 0x22,
	0x7a, 0xf3, 0xcc, 0x9b, 0xf7, 0xbc, 0xf3, 0x66, 0x06, 0xf9, 0x22, 0xa5, 0x31, 0x93, 0x10, 0xd3,
	0xb1, 0xe0, 0xd2, 0xc4, 0x66, 0x3a, 0xe1, 0x10, 0x4d, 0xb4, 0x32, 0x0a, 0x37, 0x44, 0x4a, 0x23,
	0x26, 0x21, 0x72, 0x98, 0xff, 0x38, 0x53, 0x99, 0xb2, 0x50, 0x5c, 0x7c, 0xb9, 0x2a, 0xff, 0xaf,
	0x82, 0xca, 0x73, 0x25, 0xd7, 0x15, 0xc2, 0xb9, 0x87, 0x1e, 0x0d, 0x20, 0x4b, 0x78, 0x26, 0xc0,
	0x70, 0xdd, 0x57, 0x39, 0x11, 0x12, 0x1f, 0xa3, 0x3a, 0xa8, 0x73, 0x4d, 0xf9, 0x70, 0xa2, 0xb4,
	0x69, 0x7a, 0x87, 0xde, 0xd1, 0x5e, 0x77, 0xff, 0x6a, 0xd6, 0xc2, 0x53, 0x92, 0x8f, 0x4f, 0xc2,
	0x35, 0x30, 0x4c, 0x90, 0x8b, 0x3e, 0x2a, 0x6d, 0xf0, 0x5b, 0xd4, 0x28, 0x31, 0x3a, 0x22, 0x52,
	0xf2, 0x71, 0xf3, 0x9e, 0xe5, 0x1e, 0x5c, 0xcd, 0x5a, 0x4f, 0x6e, 0x70, 0x4b, 0x3c, 0x4c, 0x1e,
	0xba, 0x44, 0xcf, 0xc5, 0x78, 0x1f, 0xd5, 0x98, 0x7d, 0x44, 0x73, 0xa7, 0x60, 0x26, 0x65, 0x84,
	0x7d, 0x74, 0x3f, 0xe7, 0x86, 0x30, 0x62, 0x48, 0xb3, 0x7a, 0xe8, 0x1d, 0x3d, 0x48, 0xae, 0xe3,
	0x82, 0x03, 0x5c, 0x32, 0xae, 0x9b, 0xbb, 0x16, 0x29, 0xa3, 0x93, 0xea, 0xf7, 0x9f, 0xad, 0x4a,
	0xf8, 0x0c, 0x1d, 0xdc, 0xea, 0x30, 0xe1, 0x30, 0x51, 0x12, 0x78, 0xf8, 0xdb, 0x43, 0xfe, 0x00,
	0x32, 0x97, 0x3d, 0x05, 0x50, 0x54, 0x10, 0x23, 0x94, 0xec, 0x69, 0x4e, 0x0c, 0x5f, 0x53, 0xf6,
	0xd6, 0x95, 0xf1, 0x31, 0xaa, 0x31, 0x09, 0x43, 0xc1, 0x6c, 0x7f, 0xf5, 0x8e, 0x1f, 0x5d, 0x4f,
	0xc2, 0x7a, 0x1c, 0xbd, 0x57, 0x94, 0x8c, 0xfb, 0x1f, 0x3e, 0x9d, 0xf5, 0xbb, 0xd5, 0x8b, 0x59,
	0xab, 0x92, 0xec, 0x32, 0x09, 0x67, 0x0c, 0x9f, 0x22, 0x04, 0x9a, 0x0e, 0xdd, 0xbc, 0x6c, 0x8b,
	0xf5, 0xce, 0xf3, 0x4d, 0x72, 0xcf, 0xa2, 0xee, 0x4d, 0x25, 0x7d, 0x0f, 0x34, 0x75, 0xe9, 0x42,
	0x82, 0x81, 0x59, 0x49, 0x54, 0xef, 0x2e, 0xc1, 0xc0, 0xb8, 0x74, 0x69, 0xcc, 0x4b, 0x14, 0x6e,
	0x6f, 0x7d, 0xe5, 0x50, 0x67, 0xee, 0xa1, 0x9d, 0x01, 0x64, 0xf8, 0x0b, 0x6a, 0x6c, 0x6c, 0xc9,
	0x8b, 0xe8, 0xe6, 0xfa, 0x45, 0xb7, 0x6c, 0xf6, 0x5f, 0xfd, 0xb7, 0x64, 0xf5, 0x1f, 0x3c, 0x45,
	0x4f, 0xb7, 0x4d, 0xe1, 0xf5, 0x3f, 0x54, 0xb6, 0xd4, 0xfa, 0x9d, 0xbb, 0xd7, 0xae, 0x7e, 0xdd,
	0x4d, 0x2f, 0x16, 0x81, 0x77, 0xb9, 0x08, 0xbc, 0xf9, 0x22, 0xf0, 0x7e, 0x2c, 0x83, 0xca, 0xe5,
	0x32, 0xa8, 0xfc, 0x5a, 0x06, 0x95, 0xcf, 0xef, 0x32, 0x61, 0x46, 0xe7, 0x69, 0xe1, 0x6a, 0x5c,
	0xac, 0x1a, 0x1d, 0x11, 0x21, 0xc7, 0x24, 0x8d, 0xa9, 0x82, 0x5c, 0x41, 0x1b, 0xd8, 0xd7, 0xb6,
	0x90, 0x86, 0x6b, 0x0b, 0xb4, 0x8b, 0x23, 0xfb, 0x16, 0x8b, 0x94, 0xb6, 0x37, 0x0f, 0x36, 0xad,
	0xd9, 0x7b, 0x7b, 0xf3, 0x27, 0x00, 0x00, 0xff, 0xff, 0x16, 0xba, 0x9d, 0x32, 0xcf, 0x03, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	RegisterDomain(ctx context.Context, in *MsgRegisterDomain, opts ...grpc.CallOption) (*MsgRegisterDomainResponse, error)
	DomainAssociationCreate(ctx context.Context, in *MsgDomainAssociationCreate, opts ...grpc.CallOption) (*MsgDomainAssociationCreateResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RegisterDomain(ctx context.Context, in *MsgRegisterDomain, opts ...grpc.CallOption) (*MsgRegisterDomainResponse, error) {
	out := new(MsgRegisterDomainResponse)
	err := c.cc.Invoke(ctx, "/ibc.dns.client.Msg/RegisterDomain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DomainAssociationCreate(ctx context.Context, in *MsgDomainAssociationCreate, opts ...grpc.CallOption) (*MsgDomainAssociationCreateResponse, error) {
	out := new(MsgDomainAssociationCreateResponse)
	err := c.cc.Invoke(ctx, "/ibc.dns.client.Msg/DomainAssociationCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	RegisterDomain(context.Context, *MsgRegisterDomain) (*MsgRegisterDomainResponse, error)
	DomainAssociationCreate(context.Context, *MsgDomainAssociationCreate) (*MsgDomainAssociationCreateResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) RegisterDomain(ctx context.Context, req *MsgRegisterDomain) (*MsgRegisterDomainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterDomain not implemented")
}
func (*UnimplementedMsgServer) DomainAssociationCreate(ctx context.Context, req *MsgDomainAssociationCreate) (*MsgDomainAssociationCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DomainAssociationCreate not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_RegisterDomain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterDomain)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterDomain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.dns.client.Msg/RegisterDomain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterDomain(ctx, req.(*MsgRegisterDomain))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DomainAssociationCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDomainAssociationCreate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DomainAssociationCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ibc.dns.client.Msg/DomainAssociationCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DomainAssociationCreate(ctx, req.(*MsgDomainAssociationCreate))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ibc.dns.client.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterDomain",
			Handler:    _Msg_RegisterDomain_Handler,
		},
		{
			MethodName: "DomainAssociationCreate",
			Handler:    _Msg_DomainAssociationCreate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ibc/dns/client/types.proto",
}

func (m *MsgRegisterDomain) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterDomain) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterDomain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Metadata) > 0 {
		i -= len(m.Metadata)
		copy(dAtA[i:], m.Metadata)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Metadata)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Domain) > 0 {
		i -= len(m.Domain)
		copy(dAtA[i:], m.Domain)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Domain)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SourceChannel) > 0 {
		i -= len(m.SourceChannel)
		copy(dAtA[i:], m.SourceChannel)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.SourceChannel)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.SourcePort) > 0 {
		i -= len(m.SourcePort)
		copy(dAtA[i:], m.SourcePort)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.SourcePort)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgRegisterDomainResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRegisterDomainResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRegisterDomainResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgDomainAssociationCreate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDomainAssociationCreate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDomainAssociationCreate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.DstClient.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.SrcClient.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.DnsId.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDomainAssociationCreateResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDomainAssociationCreateResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDomainAssociationCreateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgRegisterDomain) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SourcePort)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.SourceChannel)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Domain)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Metadata)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *MsgRegisterDomainResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgDomainAssociationCreate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.DnsId.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.SrcClient.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.DstClient.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *MsgDomainAssociationCreateResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgRegisterDomain) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MsgRegisterDomain: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterDomain: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourcePort", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourcePort = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceChannel", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceChannel = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Domain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Domain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metadata", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Metadata = append(m.Metadata[:0], dAtA[iNdEx:postIndex]...)
			if m.Metadata == nil {
				m.Metadata = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = append(m.Sender[:0], dAtA[iNdEx:postIndex]...)
			if m.Sender == nil {
				m.Sender = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *MsgRegisterDomainResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MsgRegisterDomainResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRegisterDomainResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *MsgDomainAssociationCreate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MsgDomainAssociationCreate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDomainAssociationCreate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = append(m.Sender[:0], dAtA[iNdEx:postIndex]...)
			if m.Sender == nil {
				m.Sender = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DnsId", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DnsId.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrcClient", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SrcClient.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DstClient", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DstClient.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *MsgDomainAssociationCreateResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MsgDomainAssociationCreateResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDomainAssociationCreateResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
