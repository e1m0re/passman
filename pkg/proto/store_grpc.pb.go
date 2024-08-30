// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: pkg/proto/store.proto

package store

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
	Store_GetItemsList_FullMethodName = "/store.Store/GetItemsList"
	Store_UploadItem_FullMethodName   = "/store.Store/UploadItem"
	Store_DownloadItem_FullMethodName = "/store.Store/DownloadItem"
)

// StoreClient is the client API for Store service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StoreClient interface {
	GetItemsList(ctx context.Context, in *GetItemsListRequest, opts ...grpc.CallOption) (*GetItemsListResponse, error)
	UploadItem(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadItemRequest, UploadItemResponse], error)
	DownloadItem(ctx context.Context, in *DownloadItemRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadItemResponse], error)
}

type storeClient struct {
	cc grpc.ClientConnInterface
}

func NewStoreClient(cc grpc.ClientConnInterface) StoreClient {
	return &storeClient{cc}
}

func (c *storeClient) GetItemsList(ctx context.Context, in *GetItemsListRequest, opts ...grpc.CallOption) (*GetItemsListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetItemsListResponse)
	err := c.cc.Invoke(ctx, Store_GetItemsList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) UploadItem(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadItemRequest, UploadItemResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Store_ServiceDesc.Streams[0], Store_UploadItem_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadItemRequest, UploadItemResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Store_UploadItemClient = grpc.ClientStreamingClient[UploadItemRequest, UploadItemResponse]

func (c *storeClient) DownloadItem(ctx context.Context, in *DownloadItemRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadItemResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Store_ServiceDesc.Streams[1], Store_DownloadItem_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadItemRequest, DownloadItemResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Store_DownloadItemClient = grpc.ServerStreamingClient[DownloadItemResponse]

// StoreServer is the server API for Store service.
// All implementations must embed UnimplementedStoreServer
// for forward compatibility.
type StoreServer interface {
	GetItemsList(context.Context, *GetItemsListRequest) (*GetItemsListResponse, error)
	UploadItem(grpc.ClientStreamingServer[UploadItemRequest, UploadItemResponse]) error
	DownloadItem(*DownloadItemRequest, grpc.ServerStreamingServer[DownloadItemResponse]) error
	mustEmbedUnimplementedStoreServer()
}

// UnimplementedStoreServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStoreServer struct{}

func (UnimplementedStoreServer) GetItemsList(context.Context, *GetItemsListRequest) (*GetItemsListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemsList not implemented")
}
func (UnimplementedStoreServer) UploadItem(grpc.ClientStreamingServer[UploadItemRequest, UploadItemResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadItem not implemented")
}
func (UnimplementedStoreServer) DownloadItem(*DownloadItemRequest, grpc.ServerStreamingServer[DownloadItemResponse]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadItem not implemented")
}
func (UnimplementedStoreServer) mustEmbedUnimplementedStoreServer() {}
func (UnimplementedStoreServer) testEmbeddedByValue()               {}

// UnsafeStoreServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StoreServer will
// result in compilation errors.
type UnsafeStoreServer interface {
	mustEmbedUnimplementedStoreServer()
}

func RegisterStoreServer(s grpc.ServiceRegistrar, srv StoreServer) {
	// If the following call pancis, it indicates UnimplementedStoreServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Store_ServiceDesc, srv)
}

func _Store_GetItemsList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemsListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).GetItemsList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Store_GetItemsList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).GetItemsList(ctx, req.(*GetItemsListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_UploadItem_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StoreServer).UploadItem(&grpc.GenericServerStream[UploadItemRequest, UploadItemResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Store_UploadItemServer = grpc.ClientStreamingServer[UploadItemRequest, UploadItemResponse]

func _Store_DownloadItem_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadItemRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StoreServer).DownloadItem(m, &grpc.GenericServerStream[DownloadItemRequest, DownloadItemResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Store_DownloadItemServer = grpc.ServerStreamingServer[DownloadItemResponse]

// Store_ServiceDesc is the grpc.ServiceDesc for Store service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Store_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "store.Store",
	HandlerType: (*StoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetItemsList",
			Handler:    _Store_GetItemsList_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadItem",
			Handler:       _Store_UploadItem_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadItem",
			Handler:       _Store_DownloadItem_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/proto/store.proto",
}
