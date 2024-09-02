package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/proto"
)

type GRPCClient interface {
	// Shutdown closes connection.
	Shutdown() error

	// SignUp registers new user on the server.
	SignUp(ctx context.Context, credentials model.Credentials) error

	// UploadItem send file to server.
	UploadItem(ctx context.Context, filePath string) error

	// DownloadItem gets file from server.
	DownloadItem(ctx context.Context, id string) error
}

type grpcClient struct {
	config      *ClientConfig
	connection  *grpc.ClientConn
	storeClient proto.StoreServiceClient
	authClient  proto.AuthServiceClient
}

// Shutdown closes connection.
func (g *grpcClient) Shutdown() error {
	return g.connection.Close()
}

// SignUp registers new user on the server.
func (g *grpcClient) SignUp(ctx context.Context, credentials model.Credentials) error {
	request := &proto.SignUpRequest{
		Username: credentials.Username,
		Password: credentials.Password,
	}

	resp, err := g.authClient.SignUp(ctx, request)
	if err != nil {
		return err
	}

	if resp.Status == proto.StatusCode(codes.Internal) {
		return errors.New(resp.Message)
	}

	return nil
}

// UploadItem send file to server.
func (g *grpcClient) UploadItem(ctx context.Context, id string) error {
	filePath := filepath.Join(g.config.WorkDir, id)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	stream, err := g.storeClient.UploadItem(ctx)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
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
		if err := stream.Send(&proto.UploadItemRequest{Id: id, Chunk: chunk}); err != nil {
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

	return nil
}

// DownloadItem gets file from server.
func (g *grpcClient) DownloadItem(ctx context.Context, id string) error {
	req := &proto.DownloadItemRequest{
		Guid: id,
	}

	stream, err := g.storeClient.DownloadItem(ctx, req)
	if err != nil {
		return err
	}

	filePath := filepath.Join(g.config.WorkDir, id)
	var downloaded int64
	var buffer bytes.Buffer

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			if err := os.WriteFile(filePath, buffer.Bytes(), 0666); err != nil {
				return err
			}
			break
		}
		if err != nil {
			buffer.Reset()
			return err
		}

		payload := res.GetPayload()
		size := len(payload)
		downloaded += int64(size)

		buffer.Write(payload)
	}

	return nil
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
		storeClient: proto.NewStoreServiceClient(conn),
		authClient:  proto.NewAuthServiceClient(conn),
	}, err
}
