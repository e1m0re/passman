package repository

import (
	"context"

	"passman/server/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=UserRepository
type UserRepository interface {
	// AddUser creates new user.
	AddUser(ctx context.Context, credentials models.Credentials) (*models.User, error)
	// FindUserByID finds and returns user instance by id.
	FindUserByID(ctx context.Context, id models.UserID) (*models.User, error)
	// FindUserByUsername finds and returns user instance by username.
	FindUserByUsername(ctx context.Context, username []byte) (*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=UsersDataRepository
type UsersDataRepository interface {
	// AddItem creates new users data item.
	AddItem(ctx context.Context, data models.UsersDataItemInfo) (*models.UsersDataItem, error)
	// FindItemByID finds and returns users data item.
	FindItemByID(ctx context.Context, id models.UsersDataItemID) (*models.UsersDataItem, error)
}
