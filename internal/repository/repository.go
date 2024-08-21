package repository

import (
	"context"

	"e1m0re/passman/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=UserRepository
type UserRepository interface {
	// Add creates new user.
	Add(ctx context.Context, credentials models.Credentials) (*models.User, error)
	// FindByID finds and returns user instance by id.
	FindByID(ctx context.Context, id models.UserID) (*models.User, error)
	// FindByUsername finds and returns user instance by username.
	FindByUsername(ctx context.Context, username []byte) (*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=UsersDataRepository
type UsersDataRepository interface {
	// Add creates new users data item.
	Add(ctx context.Context, data models.UsersDataItemInfo) (*models.UsersDataItem, error)
	// FindByID finds and returns users data item.
	FindByID(ctx context.Context, id models.UsersDataItemID) (*models.UsersDataItem, error)
}
