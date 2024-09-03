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

	os.Mkdir("./data", 0770)
	app1 := app.NewApp(&config.AppConfig{
		GRPCConfig: &grpc.ClientConfig{
			Port:     3000,
			Hostname: "192.168.10.102",
			WorkDir:  "./data",
		},
	})

	err := app1.Run(ctx)
	if err != nil {
		slog.Error("start application failed", slog.String("error", err.Error()))
		return
	}
}
