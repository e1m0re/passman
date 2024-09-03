package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/tools/encrypt"
)

func (a *app) getAddCredentialsPage() tview.Primitive {
	statusStr := tview.NewTextView().SetText("")

	data := model.CredentialItemData{
		Username: "",
		Password: "",
	}

	metadata := model.DatumMetadata{
		Title: "",
		Type:  model.DatumTypeText,
	}

	form := tview.NewForm().
		AddInputField("Title", "", 20, nil, func(text string) {
			metadata.Title = text
		}).
		AddInputField("Username", "", 20, nil, func(text string) {
			data.Username = text
		}).AddPasswordField("Password", "", 20, '*', func(text string) {
		data.Password = text
	}).AddButton("Save", func() {
		go func() {
			a.app.QueueUpdateDraw(func() {
				id := uuid.New().String()
				jsonData, _ := json.Marshal(data)
				filePath := filepath.Join(a.cfg.GRPCConfig.WorkDir, id)
				err := os.WriteFile(filePath, jsonData, 0666)
				if err != nil {
					statusStr.SetText(err.Error())
					return
				}
				ctx := context.Background()
				jsonMetadata, _ := json.Marshal(metadata)
				err = a.storeClient.UploadItem(ctx, id, string(jsonMetadata))
				if err != nil {
					statusStr.SetText(err.Error())
					return
				}

				checksum, _ := encrypt.FileMD5(filePath)
				a.store.AddItem(&model.DatumInfo{
					Metadata: string(jsonMetadata),
					File:     id,
					Checksum: checksum,
				})
				a.updateItemsListView()
				a.pages.SwitchToPage(MainPage)
			})
		}()
	}).AddButton("Cancel", func() {
		a.pages.SwitchToPage(MainPage)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 10, 1, true).
		AddItem(statusStr, 2, 1, false)

	return flex
}
