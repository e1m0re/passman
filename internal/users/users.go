package users

import (
	"context"

	"e1m0re/passman/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=Manager
type Manager interface {
	// AddUser creates new user.
	AddUser(ctx context.Context, credentials models.Credentials) (*models.User, error)
	// FindUserByID finds and returns user instance by id or nil.
	FindUserByID(ctx context.Context, id models.UserID) (*models.User, error)
	// FindUserByUsername finds and returns user instance by username or nil.
	FindUserByUsername(ctx context.Context, username []byte) (*models.User, error)
}
