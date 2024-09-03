package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (a *app) getAddFilePage() tview.Primitive {
	statusStr := tview.NewTextView().SetText("")

	data := model.FileItemData{
		File: "",
	}

	metadata := model.DatumMetadata{
		Title:    "",
		Type:     model.DatumTypeFile,
		FileType: "",
	}

	form := tview.NewForm().
		AddInputField("Title", "", 50, nil, func(text string) {
			metadata.Title = text
		}).
		AddInputField("File path", "", 50, nil, func(text string) {
			data.File = text
		}).
		AddButton("Save", func() {
			go func() {
				a.app.QueueUpdateDraw(func() {
					metadata.FileType = filepath.Ext(data.File)
					guid := uuid.New().String()

					err := copyFileContents(data.File, filepath.Join(a.cfg.GRPCConfig.WorkDir, guid))
					if err != nil {
						statusStr.SetText(err.Error())
						return
					}

					err = a.uploadFileToServer(guid, metadata)
					if err != nil {
						statusStr.SetText(err.Error())
						return
					}

					a.pages.SwitchToPage(MainPage)
				})
			}()
		}).
		AddButton("Cancel", func() {
			a.pages.SwitchToPage(MainPage)
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 0, 3, true).
		AddItem(statusStr, 0, 1, false)

	return flex
}
