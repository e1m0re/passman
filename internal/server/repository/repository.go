package repository

import (
	"context"

	"github.com/e1m0re/passman/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=UserRepository
type UserRepository interface {
	// AddUser creates new user.
	AddUser(ctx context.Context, credentials model.Credentials) (*model.User, error)
	// FindUserByUsername finds and returns user instance by username.
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=DatumRepository
type DatumRepository interface {
	// AddItem creates new users data item.
	AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error)
	// FindItemByFileName finds and returns users data item by filename.
	FindItemByFileName(ctx context.Context, fileName string) (*model.DatumItem, error)
	// FindByUser returns all data items by user ID.
	FindByUser(ctx context.Context, userID int) (*model.DatumItemsList, error)
}
