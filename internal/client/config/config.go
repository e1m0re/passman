package config

import "github.com/e1m0re/passman/internal/client/grpc"

type AppConfig struct {
	GRPCConfig *grpc.ClientConfig
}
