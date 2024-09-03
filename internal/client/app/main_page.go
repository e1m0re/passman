package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) getMainPage() tview.Primitive {
	var text = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(a) to add a new item\t(e) to edit selected item\t(d) to delete selected item\t(?) help \t(q) to quit")

	var itemsList = tview.NewList().ShowSecondaryText(false)

	var flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().AddItem(itemsList, 0, 1, true), 0, 6, true).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 113:
			a.app.Stop()
			//case 63:
			//	tui.pages.SwitchToPage(AboutPage)
			//case 97:
			//	tui.form.Clear(true)
			//	addCredentialsForm(model.Credentials{
			//		Password: "",
			//		Username: "",
			//	})
			//	tui.pages.SwitchToPage(AddCredentialsPage)
		}

		return event
	})

	return flex
}
