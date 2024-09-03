package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/e1m0re/passman/internal/client/app"
	"github.com/e1m0re/passman/internal/client/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cfg, err := config.NewConfig("passman.yml")
	if err != nil {
		slog.Error("failed configuration app", slog.String("error", err.Error()))
		return
	}
	app1 := app.NewApp(cfg)

	err = app1.Run(ctx)
	if err != nil {
		slog.Error("start application failed", slog.String("error", err.Error()))
		return
	}
}
