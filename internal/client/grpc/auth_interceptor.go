package grpc

import (
	"context"
	"log/slog"
	"time"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	commongrpc "github.com/e1m0re/passman/internal/common/grpc"
)

// AuthInterceptor is a client interceptor for authentication.
type AuthInterceptor struct {
	authClient  *AuthClient
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

//func (interceptor *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
//	err := interceptor.refreshToken()
//	if err != nil {
//		return err
//	}
//
//	go func() {
//		wait := refreshDuration
//		for {
//			time.Sleep(wait)
//			err := interceptor.refreshToken()
//			if err != nil {
//				wait = time.Second
//			} else {
//				wait = refreshDuration
//			}
//		}
//	}()
//
//	return nil
//}

//func (interceptor *AuthInterceptor) refreshToken() error {
//	accessToken, err := interceptor.authClient.Login()
//	if err != nil {
//		return err
//	}
//
//	interceptor.accessToken = accessToken
//	slog.Debug("access token refreshed", slog.String("new toke", accessToken))
//
//	return nil
//}

// NewAuthInterceptor initiates a new instance of AuthInterceptor.
func NewAuthInterceptor(authClient *AuthClient, refreshDuration time.Duration) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient: authClient,
	}

	//err := interceptor.scheduleRefreshToken(refreshDuration)
	//if err != nil {
	//	return nil, err
	//}

	return interceptor, nil
}
