package listener

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/e1m0re/passman/internal/server/file_server"
	transfer "github.com/e1m0re/passman/pkg/proto"
)

type grpcListener struct {
	server *grpc.Server
}

// Run starts GRPC listener.
func (l *grpcListener) Run(ctx context.Context) error {
	listen, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		return err
	}

	return l.server.Serve(listen)
}

// Shutdown stops GRPC listener.
func (l *grpcListener) Shutdown(ctx context.Context) error {
	l.server.GracefulStop()
	return nil
}

var _ Listener = (*grpcListener)(nil)

// NewGRPCListener initiates new instance of GRPC listener.
func NewGRPCListener() (Listener, error) {
	l := &grpcListener{
		server: grpc.NewServer(),
	}

	fileServer := file_server.NewServer()
	transfer.RegisterFileServiceServer(l.server, fileServer)

	return l, nil
}
