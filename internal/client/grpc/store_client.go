package grpc

import (
	"bytes"
	"context"
	"github.com/e1m0re/passman/internal/model"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/e1m0re/passman/proto"
)

// StoreClient is a client to call store management RPC.
type StoreClient struct {
	service proto.StoreServiceClient
	workDir string
}

// GetItemsList request items list from server.
func (client *StoreClient) GetItemsList(ctx context.Context) ([]*model.DatumInfo, error) {
	response, err := client.service.GetItemsList(ctx, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*model.DatumInfo, response.GetCount())
	for idx, datum := range response.GetItemsInfo() {
		result[idx] = &model.DatumInfo{
			TypeID:   model.DatumTypeID(datum.ItemType),
			UserID:   1,
			File:     datum.File,
			Checksum: datum.Checksum,
		}
	}
	return result, nil
}

// UploadItem send file to server.
func (client *StoreClient) UploadItem(ctx context.Context, id string) error {
	filePath := filepath.Join(client.workDir, id)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	stream, err := client.service.UploadItem(ctx)
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
func (client *StoreClient) DownloadItem(ctx context.Context, id string) error {
	req := &proto.DownloadItemRequest{
		Guid: id,
	}

	stream, err := client.service.DownloadItem(ctx, req)
	if err != nil {
		return err
	}

	filePath := filepath.Join(client.workDir, id)
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

// NewStoreClient initiates a new instance of StoreClient.
func NewStoreClient(cc *grpc.ClientConn, workDir string) *StoreClient {
	return &StoreClient{
		service: proto.NewStoreServiceClient(cc),
		workDir: workDir,
	}
}
