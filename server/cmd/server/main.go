package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"passman/server/internal/listeners/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	//httpListener := rest.NewRESTListener()
	//err := httpListener.Run(ctx)
	//if err != nil {
	//	if errors.Is(err, http.ErrServerClosed) {
	//		return
	//	}
	//
	//	panic(err)
	//}

	grpcListener := grpc.NewGRPCListener()
	err := grpcListener.Run(ctx)
	if err != nil {
		panic(err)
	}
}
