package users

import (
	"context"
	"errors"

	apperrors "e1m0re/passman/internal/errors"
	"e1m0re/passman/internal/models"
	"e1m0re/passman/internal/repository"
)

type usersManager struct {
	r repository.UserRepository
}

// AddUser creates new user.
func (um usersManager) AddUser(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	return um.r.AddUser(ctx, credentials)
}

// FindUserByID finds and returns user instance by id or nil.
func (um usersManager) FindUserByID(ctx context.Context, id models.UserID) (*models.User, error) {
	user, err := um.r.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrorEntityNotFound) {
			return nil, apperrors.ErrorUserNotFound
		}

		return nil, err
	}

	return user, nil
}

// FindUserByUsername finds and returns user instance by username or nil.
func (um usersManager) FindUserByUsername(ctx context.Context, username []byte) (*models.User, error) {
	user, err := um.r.FindUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrorEntityNotFound) {
			return nil, apperrors.ErrorUserNotFound
		}

		return nil, err
	}

	return user, nil
}

var _ Manager = (*usersManager)(nil)

// NewUsersManager initiates new instance of NewUsersManager.
func NewUsersManager(r repository.UserRepository) Manager {
	return &usersManager{
		r: r,
	}
}
