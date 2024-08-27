package main

import (
	"context"
	"log"
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
		log.Fatal(err)
	}

	err = grpcClient.SendFile(ctx, "/Users/elmore/Downloads/ideaIU-2024.2.0.2-aarch64.dmg")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("file send successfully")
}
