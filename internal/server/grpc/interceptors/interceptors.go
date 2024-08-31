package interceptors

import (
	"context"
	"log/slog"

	googleGrpc "google.golang.org/grpc"
)

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *googleGrpc.UnaryServerInfo,
	handler googleGrpc.UnaryHandler,
) (interface{}, error) {
	slog.Info("->> unary interceptor: ", slog.String("method", info.FullMethod))
	return handler(ctx, req)
}

func StreamInterceptor(
	srv interface{},
	stream googleGrpc.ServerStream,
	info *googleGrpc.StreamServerInfo,
	handler googleGrpc.StreamHandler,
) error {
	slog.Info("->> stream interceptor: ", slog.String("method", info.FullMethod))
	return handler(srv, stream)
}
