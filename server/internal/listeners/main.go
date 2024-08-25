package listeners

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=Listener
type Listener interface {
	// Run starts listener.
	Run(ctx context.Context) error
	// Shutdown stops listener.
	Shutdown(ctx context.Context) error
}
