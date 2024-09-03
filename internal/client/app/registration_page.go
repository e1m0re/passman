package app

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

func (a *app) getRegistrationForm() tview.Primitive {
	credentials := model.Credentials{}

	statusLabel := tview.NewTextView()
	statusLabel.SetText("Specify credentials and press \"Submit\"")

	form := tview.NewForm()
	form.
		AddInputField("Login: ", "", 20, func(t string, l rune) bool {
			return len(t) > 0
		}, func(text string) {
			credentials.Username = text
		}).
		AddInputField("Password: ", "", 20, func(t string, l rune) bool {
			return len(t) > 0
		}, func(text string) {
			credentials.Password = text
		}).
		AddFormItem(statusLabel).
		AddButton("Submit", func() {
			go func() {
				statusLabel.SetText("Registration...")
				a.app.QueueUpdateDraw(func() {
					err := a.authClient.SignUp(credentials)
					if err != nil {
						statusLabel.SetText(fmt.Sprintf("Registration failed: %s", err.Error()))
						return
					}

					a.pages.SwitchToPage(LoginPage)
				})
			}()
		}).
		AddButton("Cancel", func() {
			statusLabel.SetText("Specify credentials and press \"Submit\"")
			a.pages.SwitchToPage(LoginPage)
		})
	form.SetBorder(true).SetTitle("Registration")

	return form
}
