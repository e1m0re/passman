package store

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/server/repository"
	"github.com/e1m0re/passman/internal/tools/encrypt"
	"github.com/e1m0re/passman/proto"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=StoreService
type StoreManager interface {
	// AddItem creates new datum item.
	AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error)
	// GetUsersDataItemsList returns all data items by user ID.
	GetUsersDataItemsList(ctx context.Context, userID int) (*model.DatumItemsList, error)
	// SaveFile creates new file from stream.
	SaveFile(ctx context.Context, userID int, stream proto.StoreService_UploadItemServer) (os.FileInfo, error)
	// UploadFile sends file to stream.
	UploadFile(ctx context.Context, userID int, guid string, stream proto.StoreService_DownloadItemServer) error
}

type storeManger struct {
	workDir         string
	datumRepository repository.DatumRepository
}

// AddItem creates new datum item.
func (sm storeManger) AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error) {
	return sm.datumRepository.AddItem(ctx, datumInfo)
}

// GetUsersDataItemsList returns all data items by user ID.
func (sm storeManger) GetUsersDataItemsList(ctx context.Context, userID int) (*model.DatumItemsList, error) {
	return sm.datumRepository.FindByUser(ctx, userID)
}

// SaveFile creates new file from stream.
func (sm storeManger) SaveFile(ctx context.Context, userID int, stream proto.StoreService_UploadItemServer) (os.FileInfo, error) {
	var file *os.File
	var metadata string

	fileSize := uint32(0)
	fileSize1 := 0
	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		if file == nil {
			metadata = req.GetMetadata()
			file, err = os.Create(filepath.Join(sm.workDir, strconv.Itoa(userID), req.GetId()))
			if err != nil {
				return nil, fmt.Errorf("prepare file failed: %w", err)
			}

			defer file.Close()
		}

		chunk := req.GetChunk()
		fileSize += uint32(len(chunk))
		x, err := file.Write(chunk)
		if err != nil {
			return nil, fmt.Errorf("write data to file failed: %w", err)
		}

		fileSize1 += x
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("read file info failed: %w", err)
	}

	checksum, err := encrypt.FileMD5(file.Name())
	if err != nil {
		return nil, fmt.Errorf("checksum calculation failed: %w", err)
	}

	_, err = sm.datumRepository.AddItem(ctx, model.DatumInfo{
		UserID:   userID,
		Metadata: metadata,
		File:     fileInfo.Name(),
		Checksum: checksum,
	})
	if err != nil {
		return nil, fmt.Errorf("save info to DB failed: %w", err)
	}

	return file.Stat()
}

// UploadFile sends file to stream.
func (sm storeManger) UploadFile(ctx context.Context, userID int, guid string, stream proto.StoreService_DownloadItemServer) error {
	datumItem, err := sm.datumRepository.FindItemByFileName(ctx, guid)
	if err != nil {
		return err
	}

	if datumItem == nil || datumItem.UserID != userID {
		return fmt.Errorf("file not found")
	}

	filePath := filepath.Join(sm.workDir, strconv.Itoa(userID), guid)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	f, err := os.Open(filePath)
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

		if err := stream.Send(&proto.DownloadItemResponse{Payload: data}); err != nil {
			return err
		}

		totalBytesStreamed += int64(bytesRead)
	}

	return nil
}

var _ StoreManager = (*storeManger)(nil)

// NewStoreManager initiates new instance of StoreManager.
func NewStoreManager(workDir string, datumRepository repository.DatumRepository) StoreManager {
	return &storeManger{
		workDir:         workDir,
		datumRepository: datumRepository,
	}
}
