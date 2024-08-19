package auth

import (
	"context"

	"e1m0re/passman/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=Authenticator
type Authenticator interface {
	// Registration creates new user by credentials.
	Registration(ctx context.Context, credentials models.Credentials) (*models.User, error)
	// Login check credentials and create session.
	Login(ctx context.Context, credentials models.Credentials) (ok bool, err error)
}
