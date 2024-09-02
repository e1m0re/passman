package main

import (
	"log/slog"
	"time"

	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	grpcCtrl "github.com/e1m0re/passman/internal/controller/grpc"
	"github.com/e1m0re/passman/internal/repository"
	"github.com/e1m0re/passman/internal/server"
	"github.com/e1m0re/passman/internal/server/grpc"
	"github.com/e1m0re/passman/internal/service/db"
	"github.com/e1m0re/passman/internal/service/jwt"
	"github.com/e1m0re/passman/internal/service/store"
	"github.com/e1m0re/passman/internal/service/users"
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

	jwtManager := jwt.NewJWTManager("secretKey", time.Second*30)

	grpcServerCfg := &grpc.Config{
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
	}
	grpcServer, err := grpc.NewGRPCServer(grpcServerCfg, jwtManager)
	if err != nil {
		slog.Error("failed initiates GRPC server", slog.String("error", err.Error()))
		return
	}

	userRepository := repository.NewUserRepository(dbService)
	userManager := users.NewUserManager(userRepository)

	authController := grpcCtrl.NewAuthController(jwtManager, userManager)

	datumRepository := repository.NewDatumRepository(dbService)
	storeService := store.NewStoreService("/Users/elmore/passman/server", datumRepository)
	storeController := grpcCtrl.NewStoreController(storeService)

	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			proto.RegisterAuthServiceServer(server, authController)
			proto.RegisterStoreServiceServer(server, storeController)
		})

	server.AddShutdownHook(grpcServer, dbService)
}
