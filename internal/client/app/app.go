package app

import (
	"context"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/e1m0re/passman/internal/client/config"
	grpcclient "github.com/e1m0re/passman/internal/client/grpc"
)

var (
	BuildVersion = "0.0.1"
	BuildDate    = "03.09.2024"
)

type App interface {
	// Run starts client application.
	Run(ctx context.Context) error
}

type app struct {
	app   *tview.Application
	pages *tview.Pages

	cfg *config.AppConfig

	authInterceptor *grpcclient.AuthInterceptor
	authClient      *grpcclient.AuthClient
	storeClient     *grpcclient.StoreClient
}

// InitStoreClient initiates client for API store.
func (a *app) InitStoreClient(token string) error {
	a.authInterceptor = grpcclient.NewAuthInterceptor(token)
	secConnection, err := grpc.NewClient(
		a.cfg.GRPCConfig.GetServer(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(a.authInterceptor.Unary()),
		grpc.WithStreamInterceptor(a.authInterceptor.Stream()),
	)
	if err != nil {
		return err
	}

	a.storeClient = grpcclient.NewStoreClient(secConnection, a.cfg.GRPCConfig.WorkDir)

	return nil
}

// Run starts client application.
func (a *app) Run(ctx context.Context) error {

	anonConnection, err := grpc.NewClient(a.cfg.GRPCConfig.GetServer(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	a.authClient = grpcclient.NewAuthClient(anonConnection)

	if err := a.app.SetRoot(a.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}

	return nil
}

var _ App = (*app)(nil)

func (a *app) initTui() {
	a.app = tview.NewApplication()

	pages := tview.NewPages()
	pages.AddPage(LoginPage, a.getLoginForm(), true, true)
	pages.AddPage(RegistrationPage, a.getRegistrationForm(), true, false)
	pages.AddPage(MainPage, a.getMainPage(), true, false)
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
