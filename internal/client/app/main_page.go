package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

var itemDetail = tview.NewTextView()

func (a *app) getMainPage() tview.Primitive {
	itemDetail.SetBorder(true).SetTitle("Item info")

	var helpLine = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(a) to add a new item\t(q) to quit")
	helpLine.SetBorder(true)

	a.itemsListView.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		go func() {
			a.app.QueueUpdateDraw(func() {
				item := a.store.GetItemById(s2)
				filePath := filepath.Join(a.cfg.App.WorkDir, item.File)
				_, err := os.Stat(filePath)
				if errors.Is(err, os.ErrNotExist) {
					ctx := context.Background()
					err := a.storeClient.DownloadItem(ctx, item.File)
					if err != nil {
						itemDetail.SetText(fmt.Sprintf("Getting data error: %s", err.Error()))
						return
					}
				}

				metadata := &model.DatumMetadata{}
				json.Unmarshal([]byte(item.Metadata), metadata)
				if metadata.Type == model.DatumTypeFile {
					itemDetail.SetText(fmt.Sprintf("Title: %s\nFilename:%s", metadata.Title, metadata.File))
					return
				}

				file, err := os.Open(filePath)
				if err != nil {
					itemDetail.SetText(fmt.Sprintf("Parsing data error: %s", err.Error()))
					return
				}
				defer file.Close()

				jsonData, _ := io.ReadAll(file)
				switch metadata.Type {
				case model.DatumTypeText:
					data := &model.TextItemData{}
					json.Unmarshal(jsonData, data)
					itemDetail.SetText(fmt.Sprintf("Title: %s\nText:%s", metadata.Title, data.Text))
				case model.DatumTypeCredentials:
					data := &model.CredentialItemData{}
					err := json.Unmarshal(jsonData, data)
					if err != nil {
						itemDetail.SetText(fmt.Sprintf("error %s", err.Error()))
						return
					}
					itemDetail.SetText(fmt.Sprintf("Title: %s\nUsername:%s\nPassword:%s", metadata.Title, data.Username, data.Password))
				case model.DatumTypeCreditCard:
					data := &model.CreditCardItemData{}
					json.Unmarshal(jsonData, data)
					itemDetail.SetText(fmt.Sprintf("Title: %s\nNumber:%s\nOwner:%s\nPeriod:%s\nCVV:%s", metadata.Title, data.Number, data.Owner, data.Period, data.CVV))
				}
			})
		}()
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(a.itemsListView, 0, 1, true).
			AddItem(itemDetail, 0, 1, false),
			0, 6, false).
		AddItem(helpLine, 0, 1, false)

	flex.SetBorder(true).SetTitle(fmt.Sprintf("Passman %s (build by %s)", BuildVersion, BuildDate))

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 113:
			a.app.Stop()
		case 63:
			a.pages.SwitchToPage(AboutPage)
		case 97:
			a.pages.SwitchToPage(SelectNewItemTypePage)
		}

		return event
	})

	return flex
}
