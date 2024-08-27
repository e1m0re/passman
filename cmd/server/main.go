package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/e1m0re/passman/internal/server/listener"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	srv, _ := listener.NewGRPCListener()
	err := srv.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
