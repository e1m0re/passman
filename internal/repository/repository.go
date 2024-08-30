package repository

import (
	"context"
	"github.com/e1m0re/passman/internal/model"

	"github.com/e1m0re/passman/internal/models"
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

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=DatumRepository
type DatumRepository interface {
	// AddItem creates new users data item.
	AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error)
	// FindItemByFileName finds and returns users data item by filename.
	FindItemByFileName(ctx context.Context, fileName string) (*model.DatumItem, error)
}
