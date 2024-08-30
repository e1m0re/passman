package config

import "github.com/e1m0re/passman/internal/server/grpc"

type AppConfig struct {
	grpcConfig *grpc.Config
}
