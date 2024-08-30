package server

import (
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func AddShutdownHook(closers ...io.Closer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-c

	slog.Info("shutdown server started")
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			slog.Error("failed to stop closer", slog.String("error", err.Error()))
		}
	}

	slog.Info("shutdown server completed")
}
