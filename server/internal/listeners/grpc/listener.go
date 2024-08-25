package grpc

import (
	"context"
	"net"
	"os"

	"google.golang.org/grpc"

	"passman/server/internal/listeners"
)

type listener struct {
	gRPCServer *grpc.Server
}

// Run starts listener.
func (l listener) Run(ctx context.Context) error {
	listen, err := net.Listen("tcp", "0.0.0.0:"+os.Getenv("GRPC_SERVER_PORT"))
	if err != nil {
		return err
	}

	return l.gRPCServer.Serve(listen)
}

// Shutdown stops listener.
func (l listener) Shutdown(ctx context.Context) error {
	l.gRPCServer.GracefulStop()
	return nil
}

var _ listeners.Listener = (*listener)(nil)

// NewGRPCListener initiates new instance of GRPC listener.
func NewGRPCListener() listeners.Listener {
	return &listener{
		gRPCServer: grpc.NewServer(),
	}
}
