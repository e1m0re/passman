package users

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/server/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=UserProvider
type UserProvider interface {
	// CreateUser creates new user.
	CreateUser(ctx context.Context, credentials model.Credentials) (*model.User, error)
	// FindUserByUsername finds user by username.
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
	// CheckPassword validate specified users password.
	CheckPassword(ctx context.Context, user model.User, password string) (ok bool, err error)
}

type userProvider struct {
	userRepository repository.UserRepository
}

// CreateUser creates new user.
func (up userProvider) CreateUser(ctx context.Context, credentials model.Credentials) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	credentials.Password = string(hashedPassword)

	return up.userRepository.AddUser(ctx, credentials)
}

// FindUserByUsername finds user by username.
func (up userProvider) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return up.userRepository.FindUserByUsername(ctx, username)
}

// CheckPassword validate specified users password.
func (up userProvider) CheckPassword(ctx context.Context, user model.User, password string) (ok bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err
	}

	return true, err
}

var _ UserProvider = (*userProvider)(nil)

// NewUserProvider initiates new instance of UserProvider.
func NewUserProvider(userRepository repository.UserRepository) UserProvider {
	return &userProvider{
		userRepository: userRepository,
	}
}
