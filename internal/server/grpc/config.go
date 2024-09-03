package grpc

import "google.golang.org/grpc/keepalive"

type Config struct {
	Port            uint32
	KeepaliveParams keepalive.ServerParameters
	KeepalivePolicy keepalive.EnforcementPolicy
}
