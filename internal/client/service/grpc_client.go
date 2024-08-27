package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	transfer "github.com/e1m0re/passman/pkg/proto"
)

type GRPCClient interface {
	// Shutdown closes connection.
	Shutdown() error

	// SendFile send file to server.
	SendFile(ctx context.Context, filePath string) error
}

type grpcClient struct {
	conn         *grpc.ClientConn
	fileTransfer transfer.FileServiceClient
}

// Shutdown closes connection.
func (g *grpcClient) Shutdown() error {
	return g.conn.Close()
}

func (g *grpcClient) upload(ctx context.Context, filePath string, cancel context.CancelFunc) error {
	stream, err := g.fileTransfer.Upload(ctx)
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
		if err := stream.Send(&transfer.FileUploadRequest{FileName: filepath.Base(filePath), Chunk: chunk}); err != nil {
			return err
		}

		batchNumber += 1
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	fmt.Printf("Sent - %v bytes - %s\n", res.GetSize(), res.GetFileName())

	cancel()
	return nil
}

// SendFile send file to server.
func (g *grpcClient) SendFile(ctx context.Context, filePath string) error {
	ctx1, cancel := context.WithCancel(ctx)
	defer cancel()

	//go func(g *grpcClient) {
	if err := g.upload(ctx1, filePath, cancel); err != nil {
		fmt.Printf(err.Error())
		cancel()
	}
	//}(g)

	return nil
}

var _ GRPCClient = (*grpcClient)(nil)

func NewGRPCClient() (GRPCClient, error) {
	conn, err := grpc.NewClient("127.0.0.1:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		conn:         conn,
		fileTransfer: transfer.NewFileServiceClient(conn),
	}, err
}
