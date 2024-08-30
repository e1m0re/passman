package grpc

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	store "github.com/e1m0re/passman/pkg/proto"
)

type storeController struct {
	store.UnimplementedStoreServer
	workDir string
}

func (s *storeController) GetItemsList(ctx context.Context, request *store.GetItemsListRequest) (*store.GetItemsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

// UploadItem contains the logic for uploading a file to the server.
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
			err = file.SetFile(req.GetId(), s.workDir)
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

// DownloadItem contains the logic for downloading a file from the server.
func (s *storeController) DownloadItem(req *store.DownloadItemRequest, stream store.Store_DownloadItemServer) error {
	id := req.GetId()
	path := filepath.Join(s.workDir, id)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	totalBytesStreamed := int64(0)
	for totalBytesStreamed < fileSize {
		data := make([]byte, 1024)
		bytesRead, err := f.Read(data)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&store.DownloadItemResponse{Payload: data}); err != nil {
			return err
		}

		totalBytesStreamed += int64(bytesRead)
	}

	return nil
}

var _ store.StoreServer = (*storeController)(nil)

// NewStoreController initiates new instance of StoreServer.
func NewStoreController(workDir string) store.StoreServer {
	return &storeController{
		workDir: workDir,
	}
}
