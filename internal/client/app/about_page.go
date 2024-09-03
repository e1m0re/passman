package app

import (
	"fmt"

	"github.com/rivo/tview"
)

func (a *app) getAboutPage() *tview.Modal {
	return tview.NewModal().
		SetText(fmt.Sprintf("Verison: %s\nBuild Date: %s", "0.1.0", "03.03.03")).
		AddButtons([]string{"Ok"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.SwitchToPage(MainPage)
		})
}
