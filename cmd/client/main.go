package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/e1m0re/passman/internal/client/app"
	"github.com/e1m0re/passman/internal/client/config"
	"github.com/e1m0re/passman/internal/client/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	app1, err := app.NewApp(&config.AppConfig{
		GRPCConfig: &grpc.ClientConfig{
			Port:     3000,
			Hostname: "localhost",
			WorkDir:  "/Users/elmore/passman/client",
		},
	})
	if err != nil {
		slog.Error("initialization failed", slog.String("error", err.Error()))
		return
	}

	err = app1.Start(ctx)
	if err != nil {
		slog.Error("start application failed", slog.String("error", err.Error()))
		return
	}
}
