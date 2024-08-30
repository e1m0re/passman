package main

import (
	"log/slog"

	googleGrpc "google.golang.org/grpc"

	grpcCtrl "github.com/e1m0re/passman/internal/controller/grpc"
	"github.com/e1m0re/passman/internal/server"
	"github.com/e1m0re/passman/internal/server/grpc"
	store "github.com/e1m0re/passman/pkg/proto"
)

func main() {
	storeController := grpcCtrl.NewStoreController()

	grpcServer, err := grpc.NewGRPCServer()
	if err != nil {
		slog.Error("failed initiates GRPC server", slog.String("error", err.Error()))
		return
	}

	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			store.RegisterStoreServer(server, storeController)
		})

	server.AddShutdownHook(grpcServer)
}
