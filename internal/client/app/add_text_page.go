package app

import (
	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

func (a *app) getAddTextPage() tview.Primitive {
	statusStr := tview.NewTextView().SetText("")

	data := model.TextItemData{
		Text: "",
	}

	metadata := model.DatumMetadata{
		Title: "",
		Type:  model.DatumTypeText,
	}

	form := tview.NewForm().
		AddInputField("Title", "", 20, nil, func(text string) {
			metadata.Title = text
		}).
		AddTextArea("Text", "", 0, 10, 0, func(text string) {
			data.Text = text
		}).
		AddButton("Save", func() {
			go func() {
				a.app.QueueUpdateDraw(func() {
					err := a.uploadItemToServer(data, metadata)
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
		AddItem(form, 30, 1, true).
		AddItem(statusStr, 2, 1, false)

	return flex
}
