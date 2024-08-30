package main

import (
	"log/slog"

	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	grpcCtrl "github.com/e1m0re/passman/internal/controller/grpc"
	"github.com/e1m0re/passman/internal/server"
	"github.com/e1m0re/passman/internal/server/grpc"
	store "github.com/e1m0re/passman/pkg/proto"
)

func main() {
	storeController := grpcCtrl.NewStoreController("/Users/elmore/passman/server")

	grpcServer, err := grpc.NewGRPCServer(&grpc.Config{
		Port: 3000,
		KeepaliveParams: keepalive.ServerParameters{
			MaxConnectionIdle:     100,
			MaxConnectionAge:      7200,
			MaxConnectionAgeGrace: 60,
			Time:                  10,
			Timeout:               3,
		},
		KeepalivePolicy: keepalive.EnforcementPolicy{
			MinTime:             10,
			PermitWithoutStream: true,
		},
	})
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
