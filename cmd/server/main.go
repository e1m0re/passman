package main

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/e1m0re/passman/internal/server"
	grpcCtrl "github.com/e1m0re/passman/internal/server/controller/grpc"
	"github.com/e1m0re/passman/internal/server/grpc"
	"github.com/e1m0re/passman/internal/server/repository"
	"github.com/e1m0re/passman/internal/server/service/db"
	"github.com/e1m0re/passman/internal/server/service/jwt"
	"github.com/e1m0re/passman/internal/server/service/store"
	"github.com/e1m0re/passman/internal/server/service/users"
	"github.com/e1m0re/passman/proto"
)

func main() {
	dbService, err := db.NewDBService(db.DatabaseConfig{
		Driver:                  "pgx",
		Url:                     os.Getenv("DATABASE_DSN"),
		ConnMaxLifetimeInMinute: 3,
		MaxOpenConnections:      10,
		MaxIdleConnections:      1,
	})
	if err != nil {
		slog.Error("database connection failed", slog.String("error", err.Error()))
		return
	}

	jwtManager := jwt.NewJWTManager("secretKey", time.Minute*30)
	port, err := strconv.Atoi(os.Getenv("GRPC_SERVER_PORT"))
	if err != nil {
		slog.Error("invalid grpc port")
		return
	}

	grpcServerCfg := &grpc.Config{
		Port: uint32(port),
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

	userRepository := repository.NewUserRepository(dbService)
	userProvider := users.NewUserProvider(userRepository)

	datumRepository := repository.NewDatumRepository(dbService)
	storeService := store.NewStoreManager(os.Getenv("STORAGE_PATH"), datumRepository)

	grpcServer, err := grpc.NewGRPCServer(grpcServerCfg, jwtManager, userProvider)
	if err != nil {
		slog.Error("failed initiates GRPC server", slog.String("error", err.Error()))
		return
	}

	authController := grpcCtrl.NewAuthController(jwtManager, userProvider)
	storeController := grpcCtrl.NewStoreController(storeService)

	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			proto.RegisterAuthServiceServer(server, authController)
			proto.RegisterStoreServiceServer(server, storeController)
		})

	server.AddShutdownHook(grpcServer, dbService)
}
