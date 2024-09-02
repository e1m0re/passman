package grpc

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient interface {
	// Shutdown closes connection.
	Shutdown() error
}

type grpcClient struct {
	config      *ClientConfig
	connection  *grpc.ClientConn
	storeClient *StoreClient
	authClient  *AuthClient
}

// Shutdown closes connection.
func (client *grpcClient) Shutdown() error {
	return client.connection.Close()
}

var _ GRPCClient = (*grpcClient)(nil)

// NewGRPCClient initiates new instance of GRPCClient.
func NewGRPCClient(cfg *ClientConfig) (GRPCClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		config:      cfg,
		connection:  conn,
		storeClient: NewStoreClient(conn, cfg.WorkDir),
		authClient:  NewAuthClient(conn, "", ""),
	}, err
}
