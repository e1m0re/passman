package app

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/e1m0re/passman/internal/client/config"
	grpcclient "github.com/e1m0re/passman/internal/client/grpc"
)

var AccessToken = ""

type App interface {
	// Start runs client application.
	Start(ctx context.Context) error
}

type app struct {
	cfg *config.AppConfig
}

// Start runs client application.
func (a app) Start(ctx context.Context) error {

	username := "user"
	password := "password"
	server := fmt.Sprintf("%s:%d", a.cfg.GRPCConfig.Hostname, a.cfg.GRPCConfig.Port)
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())

	anonConnection, err := grpc.NewClient(server, transportOption)
	if err != nil {
		return err
	}

	authClient := grpcclient.NewAuthClient(anonConnection, username, password)
	interceptor, err := grpcclient.NewAuthInterceptor(authClient, 25*time.Second)
	if err != nil {
		return fmt.Errorf("create interceptors failed: %w", err)
	}

	_, err = grpc.NewClient(
		server,
		transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		return err
	}

	//_ := grpcclient.NewStoreClient(secConnection, a.cfg.GRPCConfig.WorkDir)

	time.Sleep(60 * time.Second)

	return nil
}

var _ App = (*app)(nil)

// NewApp initiates new instance of App.
func NewApp(cfg *config.AppConfig) App {
	return &app{
		cfg: cfg,
	}
}
