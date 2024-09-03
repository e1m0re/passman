package app

import (
	"context"
	"time"

	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/e1m0re/passman/internal/client/config"
	grpcclient "github.com/e1m0re/passman/internal/client/grpc"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
)

type App interface {
	// Run starts client application.
	Run(ctx context.Context) error
}

type app struct {
	app           *tview.Application
	pages         *tview.Pages
	itemsListView *tview.List

	cfg *config.Config

	authInterceptor *grpcclient.AuthInterceptor
	authClient      *grpcclient.AuthClient
	storeClient     *grpcclient.StoreClient

	store *Store
}

// InitStoreClient initiates client for API store.
func (a *app) InitStoreClient(ctx context.Context, token string) error {
	a.authInterceptor = grpcclient.NewAuthInterceptor(token)
	secConnection, err := grpc.NewClient(
		a.cfg.GetServer(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(a.authInterceptor.Unary()),
		grpc.WithStreamInterceptor(a.authInterceptor.Stream()),
	)
	if err != nil {
		return err
	}

	a.storeClient = grpcclient.NewStoreClient(secConnection, a.cfg.App.WorkDir)

	go func() {
		a.syncItemsList(ctx)
		a.app.Draw()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(a.cfg.Server.SyncInterval):
				// todo обработка ошибок
				a.syncItemsList(ctx)
			}
		}
	}()

	return nil
}

// Run starts client application.
func (a *app) Run(ctx context.Context) error {

	anonConnection, err := grpc.NewClient(a.cfg.GetServer(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	a.itemsListView = tview.NewList()
	a.itemsListView.SetBorder(true).SetTitle("List items")

	pages := tview.NewPages()
	pages.AddPage(LoginPage, a.getLoginForm(), true, true)
	pages.AddPage(RegistrationPage, a.getRegistrationForm(), true, false)
	pages.AddPage(MainPage, a.getMainPage(), true, false)
	pages.AddPage(SelectNewItemTypePage, a.getSelectItemTypePage(), true, false)
	pages.AddPage(AddCredentialsPage, a.getAddCredentialsPage(), true, false)
	pages.AddPage(AddSimpleTextPage, a.getAddTextPage(), true, false)
	pages.AddPage(AddCreditCardPage, a.getAddCreditCardPage(), true, false)
	pages.AddPage(AddFilePage, a.getAddFilePage(), true, false)

	a.pages = pages
}

// NewApp initiates new instance of App.
func NewApp(cfg *config.Config) App {
	app := &app{
		cfg:   cfg,
		store: NewStore(),
	}

	app.initTui()

	return app
}
