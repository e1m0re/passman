package grpc

import (
	"context"
	"log/slog"

	"github.com/e1m0re/passman/internal/server/service/store"
	"github.com/e1m0re/passman/proto"
)

type storeController struct {
	storeService store.StoreService

	proto.UnimplementedStoreServiceServer
}

func (s *storeController) GetItemsList(ctx context.Context, request *proto.GetItemsListRequest) (*proto.GetItemsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

// UploadItem contains the logic for uploading a file to the server.
func (s *storeController) UploadItem(stream proto.StoreService_UploadItemServer) error {
	fileInfo, err := s.storeService.SaveFile(stream.Context(), stream)
	if err != nil {
		slog.Error("save file to store failed", slog.String("error", err.Error()))
		return err
	}

	return stream.SendAndClose(&proto.UploadItemResponse{Id: fileInfo.Name(), Size: uint32(fileInfo.Size())})
}

// DownloadItem contains the logic for downloading a file from the server.
func (s *storeController) DownloadItem(req *proto.DownloadItemRequest, stream proto.StoreService_DownloadItemServer) error {
	err := s.storeService.UploadFile(stream.Context(), req.GetGuid(), stream)
	if err != nil {
		slog.Error("send file to client failed", slog.String("error", err.Error()))
	}

	return err
}

var _ proto.StoreServiceServer = (*storeController)(nil)

// NewStoreController initiates a new instance of StoreServiceServer.
func NewStoreController(storeService store.StoreService) proto.StoreServiceServer {
	return &storeController{
		storeService: storeService,
	}
}
