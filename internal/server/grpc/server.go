package grpc

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/e1m0re/passman/internal/server/grpc/interceptors"
	"github.com/e1m0re/passman/internal/server/service/jwt"
)

type Server interface {
	// Start runs GRPC server.
	Start(serviceRegister func(server *grpc.Server))
	io.Closer
}

type server struct {
	config     *Config
	grpcServer *grpc.Server
}

// Start runs Server.
func (s server) Start(serviceRegister func(server *grpc.Server)) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		slog.Error("failed to start GRPC server", slog.String("error", err.Error()))
		return
	}

	serviceRegister(s.grpcServer)

	if err := s.grpcServer.Serve(listener); err != nil {
		slog.Error("failed to start GRPC server", slog.String("error", err.Error()))
		return
	}

	slog.Info("start GRPC server success", slog.String("address", listener.Addr().String()))
}

// Close shutdowns GRPC server.
func (s server) Close() error {
	slog.Info("stop GRPC server")
	s.grpcServer.GracefulStop()
	return nil
}

var _ Server = (*server)(nil)

func buildKeepalivePolicy(config keepalive.EnforcementPolicy) keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             config.MinTime * time.Second,
		PermitWithoutStream: config.PermitWithoutStream,
	}
}

func buildKeepaliveParams(config keepalive.ServerParameters) keepalive.ServerParameters {
	return keepalive.ServerParameters{
		MaxConnectionIdle:     config.MaxConnectionIdle * time.Second,
		MaxConnectionAge:      config.MaxConnectionAge * time.Second,
		MaxConnectionAgeGrace: config.MaxConnectionAgeGrace * time.Second,
		Time:                  config.Time * time.Second,
		Timeout:               config.Timeout * time.Second,
	}
}

func buildOptions(config Config, jwtManager jwt.JWTManager) ([]grpc.ServerOption, error) {
	interceptor := interceptors.NewAuthInterceptor(jwtManager)

	return []grpc.ServerOption{
		grpc.KeepaliveParams(buildKeepaliveParams(config.KeepaliveParams)),
		grpc.KeepaliveEnforcementPolicy(buildKeepalivePolicy(config.KeepalivePolicy)),
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}, nil
}

// NewGRPCServer initiates new instance of Server.
func NewGRPCServer(cfg *Config, jwtManager jwt.JWTManager) (Server, error) {
	options, err := buildOptions(*cfg, jwtManager)
	if err != nil {
		return nil, err
	}

	return &server{
		config:     cfg,
		grpcServer: grpc.NewServer(options...),
	}, nil
}
