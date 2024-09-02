package main

import (
	"log/slog"

	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	grpcCtrl "github.com/e1m0re/passman/internal/controller/grpc"
	"github.com/e1m0re/passman/internal/repository"
	"github.com/e1m0re/passman/internal/server"
	"github.com/e1m0re/passman/internal/server/grpc"
	"github.com/e1m0re/passman/internal/service/db"
	"github.com/e1m0re/passman/internal/service/store"
	"github.com/e1m0re/passman/proto"
)

func main() {
	dbService, err := db.NewDBService(db.DatabaseConfig{
		Driver: "pgx",
		//Url:                     os.Getenv("DATABASE_DSN"),
		Url:                     "postgresql://passman:passman@127.0.0.1:5432/passman?sslmode=disable",
		ConnMaxLifetimeInMinute: 3,
		MaxOpenConnections:      10,
		MaxIdleConnections:      1,
	})
	if err != nil {
		slog.Error("database connection failed", slog.String("error", err.Error()))
		return
	}

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

	datumRepository := repository.NewDatumRepository(dbService)
	storeService := store.NewStoreService("/Users/elmore/passman/server", datumRepository)
	storeController := grpcCtrl.NewStoreController(storeService)

	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			proto.RegisterStoreServiceServer(server, storeController)
		})

	server.AddShutdownHook(grpcServer, dbService)
}
