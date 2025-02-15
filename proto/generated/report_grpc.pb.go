// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/report.proto

package generated

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
	ReportService_GetReportPrice_FullMethodName   = "/report_proto.ReportService/GetReportPrice"
	ReportService_GetReportHarvest_FullMethodName = "/report_proto.ReportService/GetReportHarvest"
)

// ReportServiceClient is the client API for ReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReportServiceClient interface {
	GetReportPrice(ctx context.Context, in *PriceParams, opts ...grpc.CallOption) (*ReportResponse, error)
	GetReportHarvest(ctx context.Context, in *HarvestParams, opts ...grpc.CallOption) (*ReportResponse, error)
}

type reportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReportServiceClient(cc grpc.ClientConnInterface) ReportServiceClient {
	return &reportServiceClient{cc}
}

func (c *reportServiceClient) GetReportPrice(ctx context.Context, in *PriceParams, opts ...grpc.CallOption) (*ReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, ReportService_GetReportPrice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportServiceClient) GetReportHarvest(ctx context.Context, in *HarvestParams, opts ...grpc.CallOption) (*ReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReportResponse)
	err := c.cc.Invoke(ctx, ReportService_GetReportHarvest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReportServiceServer is the server API for ReportService service.
// All implementations must embed UnimplementedReportServiceServer
// for forward compatibility.
type ReportServiceServer interface {
	GetReportPrice(context.Context, *PriceParams) (*ReportResponse, error)
	GetReportHarvest(context.Context, *HarvestParams) (*ReportResponse, error)
	mustEmbedUnimplementedReportServiceServer()
}

// UnimplementedReportServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedReportServiceServer struct{}

func (UnimplementedReportServiceServer) GetReportPrice(context.Context, *PriceParams) (*ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReportPrice not implemented")
}
func (UnimplementedReportServiceServer) GetReportHarvest(context.Context, *HarvestParams) (*ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReportHarvest not implemented")
}
func (UnimplementedReportServiceServer) mustEmbedUnimplementedReportServiceServer() {}
func (UnimplementedReportServiceServer) testEmbeddedByValue()                       {}

// UnsafeReportServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReportServiceServer will
// result in compilation errors.
type UnsafeReportServiceServer interface {
	mustEmbedUnimplementedReportServiceServer()
}

func RegisterReportServiceServer(s grpc.ServiceRegistrar, srv ReportServiceServer) {
	// If the following call pancis, it indicates UnimplementedReportServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ReportService_ServiceDesc, srv)
}

func _ReportService_GetReportPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PriceParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetReportPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReportService_GetReportPrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetReportPrice(ctx, req.(*PriceParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReportService_GetReportHarvest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HarvestParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetReportHarvest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReportService_GetReportHarvest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetReportHarvest(ctx, req.(*HarvestParams))
	}
	return interceptor(ctx, in, info, handler)
}

// ReportService_ServiceDesc is the grpc.ServiceDesc for ReportService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReportService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "report_proto.ReportService",
	HandlerType: (*ReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReportPrice",
			Handler:    _ReportService_GetReportPrice_Handler,
		},
		{
			MethodName: "GetReportHarvest",
			Handler:    _ReportService_GetReportHarvest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/report.proto",
}
