package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/e1m0re/passman/internal/client/service"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	grpcClient, err := service.NewGRPCClient()
	if err != nil {
		log.Panic(err)
	}

	err = grpcClient.SendFile(ctx, "/Users/elmore/Downloads/artifacts.zip")
	if err != nil {
		slog.WarnContext(ctx, "sync item failed", slog.String("error", err.Error()))
		return
	}

	log.Println("file send successfully")
}
