package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	commongrpc "github.com/e1m0re/passman/internal/common/grpc"
	"github.com/e1m0re/passman/internal/server/service/jwt"
)

// AuthInterceptor is a server interceptor for authentication and authorization.
type AuthInterceptor struct {
	jwtManager jwt.JWTManager
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authrization token is not provided")
	}

	accessToken := values[0]
	_, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return nil
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC.
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		slog.InfoContext(ctx, "--> unary auth interceptor", slog.String("method", info.FullMethod))

		if !commongrpc.AnonymousMethods[info.FullMethod] {
			err = interceptor.authorize(ctx, info.FullMethod)
			if err != nil {
				return nil, err
			}
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC.
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		slog.InfoContext(stream.Context(), "--> stream auth interceptor", slog.String("method", info.FullMethod))

		err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

// NewAuthInterceptor initiates new instance of AuthInterceptor.
func NewAuthInterceptor(jwtManager jwt.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager: jwtManager,
	}
}
