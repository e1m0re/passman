package file_server

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	transfer "github.com/e1m0re/passman/pkg/proto"
)

type FileServicesServer struct {
	transfer.UnimplementedFileServiceServer
}

func (g *FileServicesServer) Upload(stream transfer.FileService_UploadServer) error {
	file := NewFile()
	fileSize := 0

	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile(req.GetFileName(), "/Users/elmore/tmp")
		}

		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		chunk := req.GetChunk()
		fileSize += len(chunk)
		if err := file.Write(chunk); err != nil {
			return err
		}
	}

	fileName := filepath.Base(file.FilePath)
	fmt.Printf("saved file: %s, size: %d", file.FilePath, fileSize)
	return stream.SendAndClose(&transfer.FileUploadResponse{FileName: fileName, Size: uint32(fileSize)})
}

func NewServer() *FileServicesServer {
	return &FileServicesServer{}
}
