package app

import (
	"fmt"
	"log/slog"

	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/model"
)

func (a *app) getLoginForm() tview.Primitive {
	credentials := model.Credentials{}

	statusLabel := tview.NewTextView()
	statusLabel.SetText("Specify credentials and press \"Login\"")

	form := tview.NewForm().
		AddInputField("Login: ", credentials.Username, 20, func(t string, l rune) bool {
			return len(t) > 0
		}, func(text string) {
			credentials.Username = text
		}).
		AddPasswordField("Password: ", credentials.Password, 20, '*', func(text string) {
			credentials.Password = text
		}).
		AddFormItem(statusLabel).
		AddButton("Login", func() {
			go func() {
				statusLabel.SetText("Login...")
				a.app.QueueUpdateDraw(func() {
					token, err := a.authClient.Login(credentials)
					if err != nil {
						statusLabel.SetText(fmt.Sprintf("Login failed: %s", err.Error()))
						return
					}

					slog.Info("token", slog.String("a", token))
					a.pages.SwitchToPage(MainPage)
				})
			}()
		}).
		AddButton("Register", func() {
			statusLabel.SetText("Specify credentials and press \"Login\"")
			a.pages.SwitchToPage(RegistrationPage)
		})
	form.SetBorder(true).SetTitle("Login")

	return form
}
