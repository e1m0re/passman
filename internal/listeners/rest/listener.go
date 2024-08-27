package rest

import (
	"context"
	"errors"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"

	"passman/server/internal/listeners"
	"passman/server/internal/rest"
)

type listener struct {
	server *http.Server
}

// Run starts listener.
func (l *listener) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		err := l.server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	})

	grp.Go(func() error {
		<-ctx.Done()

		return l.Shutdown(ctx)
	})

	return grp.Wait()
}

// Shutdown stops listener.
func (l listener) Shutdown(ctx context.Context) error {
	return l.server.Shutdown(ctx)
}

var _ listeners.Listener = (*listener)(nil)

// NewRESTListener initiates new instance of HTTP listener.
func NewRESTListener() listeners.Listener {

	handler := rest.NewHandler()

	return &listener{
		server: &http.Server{
			Addr:    "0.0.0.0:" + os.Getenv("HTTP_SERVER_PORT"),
			Handler: handler.NewRouter(),
		},
	}
}
