package interceptors

import (
	"context"
	"log/slog"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	commongrpc "github.com/e1m0re/passman/internal/common/grpc"
	"github.com/e1m0re/passman/internal/server/service/jwt"
	"github.com/e1m0re/passman/internal/server/service/users"
)

// AuthInterceptor is a server interceptor for authentication and authorization.
type AuthInterceptor struct {
	jwtManager   jwt.JWTManager
	userProvider users.UserProvider
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authrization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	user, err := interceptor.userProvider.FindUserByUsername(ctx, claims.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "user not found")
	}

	return context.WithValue(ctx, commongrpc.UserIDMarker, user.ID), nil
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC.
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		slog.DebugContext(ctx, "--> unary auth interceptor", slog.String("method", info.FullMethod))

		newCtx := ctx
		if !commongrpc.AnonymousMethods[info.FullMethod] {
			newCtx, err = interceptor.authorize(ctx)
			if err != nil {
				return nil, err
			}
		}

		return handler(newCtx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC.
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		slog.DebugContext(stream.Context(), "--> stream auth interceptor", slog.String("method", info.FullMethod))

		newCtx, err := interceptor.authorize(stream.Context())
		if err != nil {
			return err
		}

		wrapped := middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		return handler(srv, wrapped)
	}
}

// NewAuthInterceptor initiates new instance of AuthInterceptor.
func NewAuthInterceptor(jwtManager jwt.JWTManager, userProvider users.UserProvider) *AuthInterceptor {
	return &AuthInterceptor{
		jwtManager:   jwtManager,
		userProvider: userProvider,
	}
}
