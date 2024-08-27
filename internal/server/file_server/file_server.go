package file_server

import (
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
	fileSize := uint32(0)

	defer func() {
		if err := file.OutputFile.Close(); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	for {
		req, err := stream.Recv()
		if file.FilePath == "" {
			file.SetFile(req.GetFileName(), "/tmp/")
		}

		if err == io.EOF {
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
	return stream.SendAndClose(&transfer.FileUploadResponse{FileName: fileName, Size: fileSize})
}

func NewServer() *FileServicesServer {
	return &FileServicesServer{}
}
