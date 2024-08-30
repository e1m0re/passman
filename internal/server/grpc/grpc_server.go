package grpc

import (
	"fmt"
	"io"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer interface {
	// Start runs GRPC server.
	Start(serviceRegister func(server *grpc.Server))
	io.Closer
}

type gRPCServer struct {
	grpcServer *grpc.Server
}

// Start runs GRPCServer.
func (g gRPCServer) Start(serviceRegister func(server *grpc.Server)) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 3000))
	if err != nil {
		slog.Error("failed to start GRPC server", slog.String("error", err.Error()))
		return
	}

	serviceRegister(g.grpcServer)

	if err := g.grpcServer.Serve(listener); err != nil {
		slog.Error("failed to start GRPC server", slog.String("error", err.Error()))
		return
	}

	slog.Info("start GRPC server success", slog.String("address", listener.Addr().String()))
}

// Close shutdowns GRPC server.
func (g gRPCServer) Close() error {
	slog.Info("stop GRPC server")
	g.grpcServer.GracefulStop()
	return nil
}

var _ GRPCServer = (*gRPCServer)(nil)

// NewGRPCServer initiates new instance of GRPCServer.
func NewGRPCServer() (GRPCServer, error) {
	return &gRPCServer{
		grpcServer: grpc.NewServer(),
	}, nil
}
