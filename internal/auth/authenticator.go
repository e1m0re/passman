package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"e1m0re/passman/internal/models"
	"e1m0re/passman/internal/users"
)

var (
	bcryptCost = 8
)

type authenticator struct {
	um users.Manager
}

// Registration creates new user by credentials.
func (a authenticator) Registration(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword(credentials.Password, bcryptCost)
	if err != nil {
		return nil, err
	}

	credentials.Password = hashPassword
	user, err := a.um.AddUser(ctx, credentials)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login check credentials and create session.
func (a authenticator) Login(ctx context.Context, credentials models.Credentials) (ok bool, err error) {
	user, err := a.um.FindUserByUsername(ctx, credentials.Username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, credentials.Password)
	if err != nil {
		return false, err
	}

	return true, err
}

var _ Authenticator = (*authenticator)(nil)

// NewAuthenticator initiates new instance of Authenticator.
func NewAuthenticator(um users.Manager) Authenticator {
	return &authenticator{
		um: um,
	}
}
