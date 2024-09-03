package app

import (
	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

func (a *app) getAddCredentialsPage() tview.Primitive {
	form := tview.NewForm()

	credentials := model.Credentials{
		Password: "",
		Username: "",
	}

	form.AddInputField("Title", "", 20, nil, func(text string) {
		credentials.Username = text
	})

	form.AddInputField("Username", "", 20, nil, func(text string) {
		credentials.Username = text
	})

	form.AddInputField("Password", credentials.Password, 20, nil, func(text string) {
		credentials.Password = text
	})

	form.AddButton("Save", func() {
		//items = append(items, model.DatumInfo{
		//	TypeID:   model.CredentialItem,
		//	UserID:   0,
		//	File:     "asd",
		//	Checksum: "asd",
		//})
		//updateItemsListView()
		a.pages.SwitchToPage(MainPage)
	})

	return form
}
