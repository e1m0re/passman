package app

import (
	"context"
	"log/slog"

	"github.com/e1m0re/passman/internal/client/config"
	"github.com/e1m0re/passman/internal/client/grpc"
)

type App interface {
	// Start runs client application.
	Start(ctx context.Context) error
}

type app struct {
	cfg        *config.AppConfig
	grpcClient grpc.GRPCClient
}

// Start runs client application.
func (a app) Start(ctx context.Context) error {
	//uid := "7ff351d0-5594-45b2-825e-a067e3ef242d"
	uid := "[69-07 KSC13] Руководство по эксплуатации.pdf"
	err := a.grpcClient.UploadItem(ctx, uid)
	if err != nil {
		slog.WarnContext(ctx, "sync item failed (to server)", slog.String("error", err.Error()))
	}
	slog.Info("sync item to server success", slog.String("id", uid))

	//uid = "Структура данных описания типа события.pdf"
	////uid = "7ff351d0-5594-45b2-825e-a067e3ef242e"
	//err = a.grpcClient.DownloadItem(ctx, uid)
	//if err != nil {
	//	slog.WarnContext(ctx, "sync item failed (from)", slog.String("error", err.Error()))
	//}
	//slog.Info("sync item to server success", slog.String("id", uid))

	return nil
}

var _ App = (*app)(nil)

// NewApp initiates new instance of App.
func NewApp(cfg *config.AppConfig) (App, error) {
	grpcClient, err := grpc.NewGRPCClient(cfg.GRPCConfig)
	if err != nil {
		return nil, err
	}

	return &app{
		cfg:        cfg,
		grpcClient: grpcClient,
	}, nil
}
