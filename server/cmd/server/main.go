package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"passman/server/internal/listeners/rest"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	httpListener := rest.NewListener()
	err := httpListener.Run(ctx)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}

		panic(err)
	}
}
