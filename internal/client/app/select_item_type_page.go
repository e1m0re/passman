package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var types = []string{"Credentials", "Text", "File", "Credit card"}

func (a *app) getSelectItemTypePage() tview.Primitive {
	help := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(q) to close dialog")

	typesList := tview.NewList().ShowSecondaryText(false)
	for idx, item := range types {
		typesList.AddItem(item, "", rune(49+idx), nil)
	}
	typesList.SetBorder(true).SetTitle("Select type of new item")
	typesList.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		switch s {
		case types[0]:
			a.pages.SwitchToPage(AddCredentialsPage)
		case types[1]:
			a.pages.SwitchToPage(AddSimpleTextPage)
		}
	})

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(typesList, 10, 1, true).
			AddItem(help, 0, 1, false), 30, 1, true).
		AddItem(nil, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 113:
			a.pages.SwitchToPage(MainPage)
		}

		return event
	})

	return flex
}
