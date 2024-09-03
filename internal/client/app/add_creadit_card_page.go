package app

import (
	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

func (a *app) getAddCreditCardPage() tview.Primitive {
	statusStr := tview.NewTextView().SetText("")

	data := model.CreditCardItemData{
		Number: "",
		Owner:  "",
		Period: "",
		CVV:    "",
	}

	metadata := model.DatumMetadata{
		Title: "",
		Type:  model.DatumTypeCreditCard,
	}

	form := tview.NewForm().
		AddInputField("Title", "", 20, nil, func(text string) {
			metadata.Title = text
		}).
		AddInputField("Number", "", 20, nil, func(text string) {
			data.Number = text
		}).
		AddInputField("Owner", "", 20, nil, func(text string) {
			data.Number = text
		}).
		AddInputField("Period", "", 10, nil, func(text string) {
			data.Number = text
		}).
		AddInputField("CVV", "", 10, nil, func(text string) {
			data.Number = text
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
		AddItem(form, 0, 3, true).
		AddItem(statusStr, 0, 1, false)

	return flex
}
