package users

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=UserManager
type UserManager interface {
	// CreateUser creates new user.
	CreateUser(ctx context.Context, credentials model.Credentials) (*model.User, error)
	// FindUserByUsername finds user by username.
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type userManager struct {
	userRepository repository.UserRepository
}

// CreateUser creates new user.
func (um userManager) CreateUser(ctx context.Context, credentials model.Credentials) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	credentials.Password = string(hashedPassword)

	return um.userRepository.AddUser(ctx, credentials)
}

// FindUserByUsername finds user by username.
func (um userManager) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return um.userRepository.FindUserByUsername(ctx, username)
}

var _ UserManager = (*userManager)(nil)

// NewUserManager initiates new instance of UserManager.
func NewUserManager(userRepository repository.UserRepository) UserManager {
	return &userManager{
		userRepository: userRepository,
	}
}
