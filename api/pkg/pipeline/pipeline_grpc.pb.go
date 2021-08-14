// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pipeline

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
	CreatePipeline(ctx context.Context, in *CreatePipelineRequest, opts ...grpc.CallOption) (*Pipeline, error)
	QueryPipeline(ctx context.Context, in *QueryPipelineRequest, opts ...grpc.CallOption) (*PipelineSet, error)
	DescribePipeline(ctx context.Context, in *DescribePipelineRequest, opts ...grpc.CallOption) (*Pipeline, error)
	WatchPipeline(ctx context.Context, in *WatchPipelineRequest, opts ...grpc.CallOption) (Service_WatchPipelineClient, error)
	CreateStep(ctx context.Context, in *CreateStepRequest, opts ...grpc.CallOption) (*Step, error)
	QueryStep(ctx context.Context, in *QueryStepRequest, opts ...grpc.CallOption) (*StepSet, error)
	DescribeStep(ctx context.Context, in *DescribeStepRequest, opts ...grpc.CallOption) (*Step, error)
	DeleteStep(ctx context.Context, in *DeleteStepRequest, opts ...grpc.CallOption) (*Step, error)
	CancelStep(ctx context.Context, in *CancelStepRequest, opts ...grpc.CallOption) (*Step, error)
	AuditStep(ctx context.Context, in *AuditStepRequest, opts ...grpc.CallOption) (*Step, error)
	DeletePipeline(ctx context.Context, in *DeletePipelineRequest, opts ...grpc.CallOption) (*Pipeline, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) CreatePipeline(ctx context.Context, in *CreatePipelineRequest, opts ...grpc.CallOption) (*Pipeline, error) {
	out := new(Pipeline)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/CreatePipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) QueryPipeline(ctx context.Context, in *QueryPipelineRequest, opts ...grpc.CallOption) (*PipelineSet, error) {
	out := new(PipelineSet)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/QueryPipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DescribePipeline(ctx context.Context, in *DescribePipelineRequest, opts ...grpc.CallOption) (*Pipeline, error) {
	out := new(Pipeline)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/DescribePipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) WatchPipeline(ctx context.Context, in *WatchPipelineRequest, opts ...grpc.CallOption) (Service_WatchPipelineClient, error) {
	stream, err := c.cc.NewStream(ctx, &Service_ServiceDesc.Streams[0], "/workflow.pipeline.Service/WatchPipeline", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceWatchPipelineClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Service_WatchPipelineClient interface {
	Recv() (*Pipeline, error)
	grpc.ClientStream
}

type serviceWatchPipelineClient struct {
	grpc.ClientStream
}

func (x *serviceWatchPipelineClient) Recv() (*Pipeline, error) {
	m := new(Pipeline)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *serviceClient) CreateStep(ctx context.Context, in *CreateStepRequest, opts ...grpc.CallOption) (*Step, error) {
	out := new(Step)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/CreateStep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) QueryStep(ctx context.Context, in *QueryStepRequest, opts ...grpc.CallOption) (*StepSet, error) {
	out := new(StepSet)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/QueryStep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DescribeStep(ctx context.Context, in *DescribeStepRequest, opts ...grpc.CallOption) (*Step, error) {
	out := new(Step)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/DescribeStep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteStep(ctx context.Context, in *DeleteStepRequest, opts ...grpc.CallOption) (*Step, error) {
	out := new(Step)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/DeleteStep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) CancelStep(ctx context.Context, in *CancelStepRequest, opts ...grpc.CallOption) (*Step, error) {
	out := new(Step)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/CancelStep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) AuditStep(ctx context.Context, in *AuditStepRequest, opts ...grpc.CallOption) (*Step, error) {
	out := new(Step)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/AuditStep", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeletePipeline(ctx context.Context, in *DeletePipelineRequest, opts ...grpc.CallOption) (*Pipeline, error) {
	out := new(Pipeline)
	err := c.cc.Invoke(ctx, "/workflow.pipeline.Service/DeletePipeline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	CreatePipeline(context.Context, *CreatePipelineRequest) (*Pipeline, error)
	QueryPipeline(context.Context, *QueryPipelineRequest) (*PipelineSet, error)
	DescribePipeline(context.Context, *DescribePipelineRequest) (*Pipeline, error)
	WatchPipeline(*WatchPipelineRequest, Service_WatchPipelineServer) error
	CreateStep(context.Context, *CreateStepRequest) (*Step, error)
	QueryStep(context.Context, *QueryStepRequest) (*StepSet, error)
	DescribeStep(context.Context, *DescribeStepRequest) (*Step, error)
	DeleteStep(context.Context, *DeleteStepRequest) (*Step, error)
	CancelStep(context.Context, *CancelStepRequest) (*Step, error)
	AuditStep(context.Context, *AuditStepRequest) (*Step, error)
	DeletePipeline(context.Context, *DeletePipelineRequest) (*Pipeline, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) CreatePipeline(context.Context, *CreatePipelineRequest) (*Pipeline, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePipeline not implemented")
}
func (UnimplementedServiceServer) QueryPipeline(context.Context, *QueryPipelineRequest) (*PipelineSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPipeline not implemented")
}
func (UnimplementedServiceServer) DescribePipeline(context.Context, *DescribePipelineRequest) (*Pipeline, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribePipeline not implemented")
}
func (UnimplementedServiceServer) WatchPipeline(*WatchPipelineRequest, Service_WatchPipelineServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchPipeline not implemented")
}
func (UnimplementedServiceServer) CreateStep(context.Context, *CreateStepRequest) (*Step, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStep not implemented")
}
func (UnimplementedServiceServer) QueryStep(context.Context, *QueryStepRequest) (*StepSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryStep not implemented")
}
func (UnimplementedServiceServer) DescribeStep(context.Context, *DescribeStepRequest) (*Step, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeStep not implemented")
}
func (UnimplementedServiceServer) DeleteStep(context.Context, *DeleteStepRequest) (*Step, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStep not implemented")
}
func (UnimplementedServiceServer) CancelStep(context.Context, *CancelStepRequest) (*Step, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelStep not implemented")
}
func (UnimplementedServiceServer) AuditStep(context.Context, *AuditStepRequest) (*Step, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuditStep not implemented")
}
func (UnimplementedServiceServer) DeletePipeline(context.Context, *DeletePipelineRequest) (*Pipeline, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePipeline not implemented")
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

func _Service_CreatePipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreatePipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/CreatePipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreatePipeline(ctx, req.(*CreatePipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_QueryPipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).QueryPipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/QueryPipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).QueryPipeline(ctx, req.(*QueryPipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DescribePipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribePipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DescribePipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/DescribePipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DescribePipeline(ctx, req.(*DescribePipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_WatchPipeline_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchPipelineRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServiceServer).WatchPipeline(m, &serviceWatchPipelineServer{stream})
}

type Service_WatchPipelineServer interface {
	Send(*Pipeline) error
	grpc.ServerStream
}

type serviceWatchPipelineServer struct {
	grpc.ServerStream
}

func (x *serviceWatchPipelineServer) Send(m *Pipeline) error {
	return x.ServerStream.SendMsg(m)
}

func _Service_CreateStep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateStep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/CreateStep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateStep(ctx, req.(*CreateStepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_QueryStep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryStepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).QueryStep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/QueryStep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).QueryStep(ctx, req.(*QueryStepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DescribeStep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeStepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DescribeStep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/DescribeStep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DescribeStep(ctx, req.(*DescribeStepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteStep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteStep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/DeleteStep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteStep(ctx, req.(*DeleteStepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_CancelStep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelStepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CancelStep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/CancelStep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CancelStep(ctx, req.(*CancelStepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_AuditStep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuditStepRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).AuditStep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/AuditStep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).AuditStep(ctx, req.(*AuditStepRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeletePipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeletePipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workflow.pipeline.Service/DeletePipeline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeletePipeline(ctx, req.(*DeletePipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "workflow.pipeline.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePipeline",
			Handler:    _Service_CreatePipeline_Handler,
		},
		{
			MethodName: "QueryPipeline",
			Handler:    _Service_QueryPipeline_Handler,
		},
		{
			MethodName: "DescribePipeline",
			Handler:    _Service_DescribePipeline_Handler,
		},
		{
			MethodName: "CreateStep",
			Handler:    _Service_CreateStep_Handler,
		},
		{
			MethodName: "QueryStep",
			Handler:    _Service_QueryStep_Handler,
		},
		{
			MethodName: "DescribeStep",
			Handler:    _Service_DescribeStep_Handler,
		},
		{
			MethodName: "DeleteStep",
			Handler:    _Service_DeleteStep_Handler,
		},
		{
			MethodName: "CancelStep",
			Handler:    _Service_CancelStep_Handler,
		},
		{
			MethodName: "AuditStep",
			Handler:    _Service_AuditStep_Handler,
		},
		{
			MethodName: "DeletePipeline",
			Handler:    _Service_DeletePipeline_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchPipeline",
			Handler:       _Service_WatchPipeline_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/pkg/pipeline/pb/pipeline.proto",
}
