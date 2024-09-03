package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/tools/encrypt"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"path/filepath"
)

var typesDescriptionMap = map[int]string{
	model.DatumTypeCredentials: "Credentials",
	model.DatumTypeText:        "Text",
	model.DatumTypeFile:        "File",
	model.DatumTypeCreditCard:  "Credit card",
}

func (a *app) updateItemsListView() {
	a.itemsListView.Clear()
	idx := 0
	for _, item := range a.store.itemsList {
		metadata := &model.DatumMetadata{}
		// Todo обработать ошибки
		err := json.Unmarshal([]byte(item.Metadata), metadata)
		if err != nil {
			s := fmt.Sprintf("%s", err.Error())
			slog.Info(s)
			continue
		}
		a.itemsListView.AddItem(metadata.Title, typesDescriptionMap[metadata.Type], rune(49+idx), nil)
		idx++
	}
}

func (a *app) uploadItemToServer(data any, metadata model.DatumMetadata) error {
	id := uuid.New().String()
	jsonData, _ := json.Marshal(data)
	filePath := filepath.Join(a.cfg.GRPCConfig.WorkDir, id)
	err := os.WriteFile(filePath, jsonData, 0666)
	if err != nil {
		return err
	}
	ctx := context.Background()
	jsonMetadata, _ := json.Marshal(metadata)
	err = a.storeClient.UploadItem(ctx, id, string(jsonMetadata))
	if err != nil {
		return err
	}

	checksum, _ := encrypt.FileMD5(filePath)
	a.store.AddItem(&model.DatumInfo{
		Metadata: string(jsonMetadata),
		File:     id,
		Checksum: checksum,
	})
	a.updateItemsListView()

	return nil
}
