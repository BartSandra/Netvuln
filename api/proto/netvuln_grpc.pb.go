// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: api/proto/netvuln.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	NetVulnService_CheckVuln_FullMethodName = "/netvuln.v1.NetVulnService/CheckVuln"
)

// NetVulnServiceClient is the client API for NetVulnService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NetVulnServiceClient interface {
	CheckVuln(ctx context.Context, in *CheckVulnRequest, opts ...grpc.CallOption) (*CheckVulnResponse, error)
}

type netVulnServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetVulnServiceClient(cc grpc.ClientConnInterface) NetVulnServiceClient {
	return &netVulnServiceClient{cc}
}

func (c *netVulnServiceClient) CheckVuln(ctx context.Context, in *CheckVulnRequest, opts ...grpc.CallOption) (*CheckVulnResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckVulnResponse)
	err := c.cc.Invoke(ctx, NetVulnService_CheckVuln_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetVulnServiceServer is the server API for NetVulnService service.
// All implementations must embed UnimplementedNetVulnServiceServer
// for forward compatibility.
type NetVulnServiceServer interface {
	CheckVuln(context.Context, *CheckVulnRequest) (*CheckVulnResponse, error)
	mustEmbedUnimplementedNetVulnServiceServer()
}

// UnimplementedNetVulnServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNetVulnServiceServer struct{}

func (UnimplementedNetVulnServiceServer) CheckVuln(context.Context, *CheckVulnRequest) (*CheckVulnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckVuln not implemented")
}
func (UnimplementedNetVulnServiceServer) mustEmbedUnimplementedNetVulnServiceServer() {}
func (UnimplementedNetVulnServiceServer) testEmbeddedByValue()                        {}

// UnsafeNetVulnServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NetVulnServiceServer will
// result in compilation errors.
type UnsafeNetVulnServiceServer interface {
	mustEmbedUnimplementedNetVulnServiceServer()
}

func RegisterNetVulnServiceServer(s grpc.ServiceRegistrar, srv NetVulnServiceServer) {
	// If the following call pancis, it indicates UnimplementedNetVulnServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NetVulnService_ServiceDesc, srv)
}

func _NetVulnService_CheckVuln_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckVulnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetVulnServiceServer).CheckVuln(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NetVulnService_CheckVuln_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetVulnServiceServer).CheckVuln(ctx, req.(*CheckVulnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NetVulnService_ServiceDesc is the grpc.ServiceDesc for NetVulnService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NetVulnService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "netvuln.v1.NetVulnService",
	HandlerType: (*NetVulnServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckVuln",
			Handler:    _NetVulnService_CheckVuln_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/netvuln.proto",
}
