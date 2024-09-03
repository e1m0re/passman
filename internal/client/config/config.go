package config

import (
	"time"

	"github.com/e1m0re/passman/internal/client/grpc"
)

type AppConfig struct {
	GRPCConfig   *grpc.ClientConfig
	SyncInterval time.Duration
}
