package service

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	store "github.com/e1m0re/passman/pkg/proto"
)

type GRPCClient interface {
	// Shutdown closes connection.
	Shutdown() error

	// SendFile send file to server.
	SendFile(ctx context.Context, filePath string) error
}

type grpcClient struct {
	conn   *grpc.ClientConn
	client store.StoreClient
}

// Shutdown closes connection.
func (g *grpcClient) Shutdown() error {
	return g.conn.Close()
}

// SendFile send file to server.
func (g *grpcClient) SendFile(ctx context.Context, filePath string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := g.client.UploadItem(ctx)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	buf := make([]byte, 1000)
	batchNumber := 1
	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunk := buf[:num]
		id := uuid.New().String()
		if err := stream.Send(&store.UploadItemRequest{Id: id, Chunk: chunk}); err != nil {
			return err
		}

		batchNumber++
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	slog.Debug(
		"sending item finished successfully",
		slog.String("file", filePath),
		slog.String("id", res.GetId()),
		slog.Int("size", int(res.GetSize())),
	)

	cancel()
	return nil
}

var _ GRPCClient = (*grpcClient)(nil)

func NewGRPCClient() (GRPCClient, error) {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		conn:   conn,
		client: store.NewStoreClient(conn),
	}, err
}
