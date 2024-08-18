package repository

import (
	"context"

	"e1m0re/passman/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.1 --name=UserRepository
type UserRepository interface {
	// AddUser creates new user.
	AddUser(ctx context.Context, userInfo models.UserInfo) (*models.User, error)
	// FindUserById finds and returns user instance by id or nil.
	FindUserById(ctx context.Context, id models.UserID) (*models.User, error)
	// FindUserByUsername finds and returns user instance by username or nil.
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
}
