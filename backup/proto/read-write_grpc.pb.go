// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: read-write.proto

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
	BackupService_BackupBlock_FullMethodName = "/backup.BackupService/BackupBlock"
)

// BackupServiceClient is the client API for BackupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BackupServiceClient interface {
	BackupBlock(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[DataRequest, DataResponse], error)
}

type backupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBackupServiceClient(cc grpc.ClientConnInterface) BackupServiceClient {
	return &backupServiceClient{cc}
}

func (c *backupServiceClient) BackupBlock(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[DataRequest, DataResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &BackupService_ServiceDesc.Streams[0], BackupService_BackupBlock_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DataRequest, DataResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BackupService_BackupBlockClient = grpc.ClientStreamingClient[DataRequest, DataResponse]

// BackupServiceServer is the server API for BackupService service.
// All implementations must embed UnimplementedBackupServiceServer
// for forward compatibility.
type BackupServiceServer interface {
	BackupBlock(grpc.ClientStreamingServer[DataRequest, DataResponse]) error
	mustEmbedUnimplementedBackupServiceServer()
}

// UnimplementedBackupServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBackupServiceServer struct{}

func (UnimplementedBackupServiceServer) BackupBlock(grpc.ClientStreamingServer[DataRequest, DataResponse]) error {
	return status.Errorf(codes.Unimplemented, "method BackupBlock not implemented")
}
func (UnimplementedBackupServiceServer) mustEmbedUnimplementedBackupServiceServer() {}
func (UnimplementedBackupServiceServer) testEmbeddedByValue()                       {}

// UnsafeBackupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BackupServiceServer will
// result in compilation errors.
type UnsafeBackupServiceServer interface {
	mustEmbedUnimplementedBackupServiceServer()
}

func RegisterBackupServiceServer(s grpc.ServiceRegistrar, srv BackupServiceServer) {
	// If the following call pancis, it indicates UnimplementedBackupServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BackupService_ServiceDesc, srv)
}

func _BackupService_BackupBlock_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BackupServiceServer).BackupBlock(&grpc.GenericServerStream[DataRequest, DataResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type BackupService_BackupBlockServer = grpc.ClientStreamingServer[DataRequest, DataResponse]

// BackupService_ServiceDesc is the grpc.ServiceDesc for BackupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BackupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "backup.BackupService",
	HandlerType: (*BackupServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BackupBlock",
			Handler:       _BackupService_BackupBlock_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "read-write.proto",
}