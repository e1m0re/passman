package grpc

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"path/filepath"

	store "github.com/e1m0re/passman/pkg/proto"
)

type storeController struct {
	store.UnimplementedStoreServer
}

func (s *storeController) GetItemsList(ctx context.Context, request *store.GetItemsListRequest) (*store.GetItemsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *storeController) UploadItem(stream store.Store_UploadItemServer) error {
	file := NewFile()
	fileSize := uint32(0)

	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			slog.Warn("close file error", slog.String("error", err.Error()))
		}
	}()

	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			err = file.SetFile(req.GetId(), "/Users/elmore/tmp")
			if err != nil {
				slog.Warn("preparing file failed", slog.String("error", err.Error()))
			}
		}

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		if err := file.Write(chunk); err != nil {
			return err
		}
	}

	fileName := filepath.Base(file.FilePath)
	slog.Info("getting data finished successfully", slog.String("file", file.FilePath), slog.Int("size", int(fileSize)))

	return stream.SendAndClose(&store.UploadItemResponse{Id: fileName, Size: fileSize})
}

var _ store.StoreServer = (*storeController)(nil)

func NewStoreController() store.StoreServer {
	return &storeController{}
}
