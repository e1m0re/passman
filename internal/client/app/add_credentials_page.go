package app

import (
	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

func (a *app) getAddCredentialsPage() tview.Primitive {
	statusStr := tview.NewTextView().SetText("")

	data := model.CredentialItemData{
		Username: "",
		Password: "",
	}

	metadata := model.DatumMetadata{
		Title: "",
		Type:  model.DatumTypeCredentials,
	}

	form := tview.NewForm().
		AddInputField("Title", "", 20, nil, func(text string) {
			metadata.Title = text
		}).
		AddInputField("Username", "", 20, nil, func(text string) {
			data.Username = text
		}).
		AddPasswordField("Password", "", 20, '*', func(text string) {
			data.Password = text
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
		AddItem(form, 10, 1, true).
		AddItem(statusStr, 2, 1, false)

	return flex
}
