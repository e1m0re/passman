package app

import (
	"context"
	"fmt"
	grpcclient "github.com/e1m0re/passman/internal/client/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rivo/tview"

	"github.com/e1m0re/passman/internal/client/config"
)

type App interface {
	// Run starts client application.
	Run(ctx context.Context) error
}

type app struct {
	app   *tview.Application
	pages *tview.Pages

	cfg *config.AppConfig

	authClient *grpcclient.AuthClient
}

// Run starts client application.
func (a *app) Run(ctx context.Context) error {

	server := fmt.Sprintf("%s:%d", a.cfg.GRPCConfig.Hostname, a.cfg.GRPCConfig.Port)
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	anonConnection, err := grpc.NewClient(server, transportOption)
	if err != nil {
		return err
	}

	a.authClient = grpcclient.NewAuthClient(anonConnection)

	if err := a.app.SetRoot(a.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}

	//username := "user"
	//password := "password"

	//interceptor, err := grpcclient.NewAuthInterceptor(authClient, 25*time.Second)
	//if err != nil {
	//	return fmt.Errorf("create interceptors failed: %w", err)
	//}
	//
	//secConnection, err := grpc.NewClient(
	//	server,
	//	transportOption,
	//	grpc.WithUnaryInterceptor(interceptor.Unary()),
	//	grpc.WithStreamInterceptor(interceptor.Stream()),
	//)
	//if err != nil {
	//	return err
	//}
	//
	//storeClient := grpcclient.NewStoreClient(secConnection, a.cfg.GRPCConfig.WorkDir)
	//dataItems, err := storeClient.GetItemsList(ctx)
	//if err != nil {
	//	slog.Error("getting data info from server failed", slog.String("error", err.Error()))
	//}
	//
	//for _, item := range dataItems {
	//	slog.Info("item", slog.String("data", fmt.Sprintf("%v", item)))
	//}
	//
	//time.Sleep(40 * time.Second)

	return nil
}

var _ App = (*app)(nil)

func (a *app) initTui() {
	a.app = tview.NewApplication()

	pages := tview.NewPages()
	pages.AddPage(LoginPage, a.getLoginForm(), true, true)
	pages.AddPage(RegistrationPage, a.getRegistrationForm(), true, false)
	pages.AddPage(MainPage, a.getMainPage(), true, false)
	pages.AddPage(AboutPage, a.getAboutPage(), true, false)
	pages.AddPage(AddCredentialsPage, a.getAddCredentialsPage(), true, false)

	a.pages = pages
}

// NewApp initiates new instance of App.
func NewApp(cfg *config.AppConfig) App {
	app := &app{
		cfg: cfg,
	}

	app.initTui()

	return app
}
