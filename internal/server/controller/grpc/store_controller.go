package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commongrpc "github.com/e1m0re/passman/internal/common/grpc"
	"github.com/e1m0re/passman/internal/server/service/store"
	"github.com/e1m0re/passman/proto"
)

type storeController struct {
	storeManager store.StoreManager

	proto.UnimplementedStoreServiceServer
}

func (sc *storeController) GetItemsList(ctx context.Context, req *proto.GetItemsListRequest) (*proto.GetItemsListResponse, error) {
	userID, ok := ctx.Value(commongrpc.UserIDMarker).(int)
	if !ok {
		return nil, status.Errorf(codes.Internal, "unknown user")
	}

	items, err := sc.storeManager.GetUsersDataItemsList(ctx, userID)
	if err != nil {
		slog.Warn("getting users data list failed", slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	count := len(*items)
	responseData := make([]*proto.ItemInfo, count)
	for idx, item := range *items {
		responseData[idx] = &proto.ItemInfo{
			Id:       item.File,
			ItemType: proto.ItemType(item.TypeID),
			Metadata: item.Metadata,
			Checksum: item.Checksum,
		}
	}
	return &proto.GetItemsListResponse{
		ItemsInfo: responseData,
		Count:     int32(count),
	}, nil
}

// UploadItem contains the logic for uploading a file to the server.
func (sc *storeController) UploadItem(stream proto.StoreService_UploadItemServer) error {
	userID, ok := stream.Context().Value(commongrpc.UserIDMarker).(int)
	if !ok {
		return status.Errorf(codes.Internal, "unknown user")
	}

	fileInfo, err := sc.storeManager.SaveFile(stream.Context(), userID, stream)
	if err != nil {
		slog.Error("save file to store failed", slog.String("error", err.Error()))
		return err
	}

	return stream.SendAndClose(&proto.UploadItemResponse{Id: fileInfo.Name(), Size: uint32(fileInfo.Size())})
}

// DownloadItem contains the logic for downloading a file from the server.
func (sc *storeController) DownloadItem(req *proto.DownloadItemRequest, stream proto.StoreService_DownloadItemServer) error {
	userID, ok := stream.Context().Value(commongrpc.UserIDMarker).(int)
	if !ok {
		return status.Errorf(codes.Internal, "unknown user")
	}

	err := sc.storeManager.UploadFile(stream.Context(), userID, req.GetGuid(), stream)
	if err != nil {
		slog.Error("send file to client failed", slog.String("error", err.Error()))
	}

	return err
}

var _ proto.StoreServiceServer = (*storeController)(nil)

// NewStoreController initiates a new instance of StoreServiceServer.
func NewStoreController(storeService store.StoreManager) proto.StoreServiceServer {
	return &storeController{
		storeManager: storeService,
	}
}
