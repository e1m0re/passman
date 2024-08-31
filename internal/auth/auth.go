package auth

import (
	"context"

	"github.com/e1m0re/passman/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=Authenticator
type Authenticator interface {
	// Registration creates new user by credentials.
	Registration(ctx context.Context, credentials model.Credentials) (*model.User, error)
	// Login check credentials and create session.
	Login(ctx context.Context, credentials model.Credentials) (ok bool, err error)
}
