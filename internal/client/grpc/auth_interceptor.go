package grpc

import (
	"context"
	"log/slog"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	commongrpc "github.com/e1m0re/passman/internal/common/grpc"
)

// AuthInterceptor is a client interceptor for authentication.
type AuthInterceptor struct {
	accessToken string
}

// Unary returns a client interceptor to authenticate unary RPC
func (interceptor *AuthInterceptor) Unary() googlegrpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *googlegrpc.ClientConn,
		invoker googlegrpc.UnaryInvoker,
		opts ...googlegrpc.CallOption,
	) error {
		slog.DebugContext(ctx, "--> unary interceptor", slog.String("method", method))

		if commongrpc.AnonymousMethods[method] {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

// Stream returns a client interceptor to authenticate stream RPC
func (interceptor *AuthInterceptor) Stream() googlegrpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *googlegrpc.StreamDesc,
		cc *googlegrpc.ClientConn,
		method string,
		streamer googlegrpc.Streamer,
		opts ...googlegrpc.CallOption,
	) (googlegrpc.ClientStream, error) {
		slog.DebugContext(ctx, "--> stream interceptor", slog.String("method", method))

		if commongrpc.AnonymousMethods[method] {
			return streamer(ctx, desc, cc, method, opts...)
		}

		return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

// NewAuthInterceptor initiates a new instance of AuthInterceptor.
func NewAuthInterceptor(token string) *AuthInterceptor {
	return &AuthInterceptor{
		accessToken: token,
	}
}
