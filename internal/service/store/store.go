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
	"github.com/e1m0re/passman/internal/repository"
	"github.com/e1m0re/passman/internal/tools/encrypt"
	proto "github.com/e1m0re/passman/proto"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=StoreService
type StoreService interface {
	// AddItem creates new datum item.
	AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error)
	// SaveFile creates new file from stream.
	SaveFile(ctx context.Context, stream proto.StoreService_UploadItemServer) (os.FileInfo, error)
	// UploadFile sends file to stream.
	UploadFile(ctx context.Context, id string, stream proto.StoreService_DownloadItemServer) error
}

type storeService struct {
	workDir         string
	datumRepository repository.DatumRepository
}

// AddItem creates new datum item.
func (s storeService) AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error) {
	return s.datumRepository.AddItem(ctx, datumInfo)
}

// SaveFile creates new file from stream.
func (s storeService) SaveFile(ctx context.Context, stream proto.StoreService_UploadItemServer) (os.FileInfo, error) {
	userId := 1
	var file *os.File

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
			file, err = os.Create(filepath.Join(s.workDir, strconv.Itoa(userId), req.GetId()))
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

	_, err = s.datumRepository.AddItem(ctx, model.DatumInfo{
		UserID:   userId,
		TypeID:   model.TextItem,
		File:     fileInfo.Name(),
		Checksum: checksum,
	})
	if err != nil {
		return nil, fmt.Errorf("save info to DB failed: %w", err)
	}

	return file.Stat()
}

// UploadFile sends file to stream.
func (s storeService) UploadFile(ctx context.Context, id string, stream proto.StoreService_DownloadItemServer) error {
	userId := 1
	datumItem, err := s.datumRepository.FindItemByFileName(ctx, id)
	if err != nil {
		return err
	}

	if datumItem == nil || datumItem.UserID != userId {
		return fmt.Errorf("file not found")
	}

	filePath := filepath.Join(s.workDir, strconv.Itoa(userId), id)
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

var _ StoreService = (*storeService)(nil)

// NewStoreService initiates new instance of StoreService.
func NewStoreService(workDir string, datumRepository repository.DatumRepository) StoreService {
	return &storeService{
		workDir:         workDir,
		datumRepository: datumRepository,
	}
}
