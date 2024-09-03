package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (a *app) getMainPage() tview.Primitive {
	var helpLine = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText("(a) to add a new item\t(q) to quit")

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.itemsListView, 0, 1, true).
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
