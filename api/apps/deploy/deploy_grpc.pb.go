// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package deploy

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	CreateApplicationDeploy(ctx context.Context, in *CreateApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploy, error)
	QueryApplicationDeploy(ctx context.Context, in *QueryApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploySet, error)
	DescribeApplicationDeploy(ctx context.Context, in *DescribeApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploy, error)
	DeleteApplicationDeploy(ctx context.Context, in *DeleteApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploy, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) CreateApplicationDeploy(ctx context.Context, in *CreateApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploy, error) {
	out := new(ApplicationDeploy)
	err := c.cc.Invoke(ctx, "/infraboard.workflow.deploy.Service/CreateApplicationDeploy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) QueryApplicationDeploy(ctx context.Context, in *QueryApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploySet, error) {
	out := new(ApplicationDeploySet)
	err := c.cc.Invoke(ctx, "/infraboard.workflow.deploy.Service/QueryApplicationDeploy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DescribeApplicationDeploy(ctx context.Context, in *DescribeApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploy, error) {
	out := new(ApplicationDeploy)
	err := c.cc.Invoke(ctx, "/infraboard.workflow.deploy.Service/DescribeApplicationDeploy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteApplicationDeploy(ctx context.Context, in *DeleteApplicationDeployRequest, opts ...grpc.CallOption) (*ApplicationDeploy, error) {
	out := new(ApplicationDeploy)
	err := c.cc.Invoke(ctx, "/infraboard.workflow.deploy.Service/DeleteApplicationDeploy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	CreateApplicationDeploy(context.Context, *CreateApplicationDeployRequest) (*ApplicationDeploy, error)
	QueryApplicationDeploy(context.Context, *QueryApplicationDeployRequest) (*ApplicationDeploySet, error)
	DescribeApplicationDeploy(context.Context, *DescribeApplicationDeployRequest) (*ApplicationDeploy, error)
	DeleteApplicationDeploy(context.Context, *DeleteApplicationDeployRequest) (*ApplicationDeploy, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) CreateApplicationDeploy(context.Context, *CreateApplicationDeployRequest) (*ApplicationDeploy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateApplicationDeploy not implemented")
}
func (UnimplementedServiceServer) QueryApplicationDeploy(context.Context, *QueryApplicationDeployRequest) (*ApplicationDeploySet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryApplicationDeploy not implemented")
}
func (UnimplementedServiceServer) DescribeApplicationDeploy(context.Context, *DescribeApplicationDeployRequest) (*ApplicationDeploy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeApplicationDeploy not implemented")
}
func (UnimplementedServiceServer) DeleteApplicationDeploy(context.Context, *DeleteApplicationDeployRequest) (*ApplicationDeploy, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteApplicationDeploy not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_CreateApplicationDeploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateApplicationDeployRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateApplicationDeploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/infraboard.workflow.deploy.Service/CreateApplicationDeploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateApplicationDeploy(ctx, req.(*CreateApplicationDeployRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_QueryApplicationDeploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryApplicationDeployRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).QueryApplicationDeploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/infraboard.workflow.deploy.Service/QueryApplicationDeploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).QueryApplicationDeploy(ctx, req.(*QueryApplicationDeployRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DescribeApplicationDeploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeApplicationDeployRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DescribeApplicationDeploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/infraboard.workflow.deploy.Service/DescribeApplicationDeploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DescribeApplicationDeploy(ctx, req.(*DescribeApplicationDeployRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteApplicationDeploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteApplicationDeployRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteApplicationDeploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/infraboard.workflow.deploy.Service/DeleteApplicationDeploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteApplicationDeploy(ctx, req.(*DeleteApplicationDeployRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "infraboard.workflow.deploy.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateApplicationDeploy",
			Handler:    _Service_CreateApplicationDeploy_Handler,
		},
		{
			MethodName: "QueryApplicationDeploy",
			Handler:    _Service_QueryApplicationDeploy_Handler,
		},
		{
			MethodName: "DescribeApplicationDeploy",
			Handler:    _Service_DescribeApplicationDeploy_Handler,
		},
		{
			MethodName: "DeleteApplicationDeploy",
			Handler:    _Service_DeleteApplicationDeploy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/app/deploy/pb/deploy.proto",
}