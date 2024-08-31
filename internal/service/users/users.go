package users

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/repository"
)

type UsersService interface {
	// NewUser creates new user.
	NewUser(ctx context.Context, credentials model.Credentials) (*model.User, error)
	// FindUserByUsername finds user by username.
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

// NewUser creates new user.
func (s userService) NewUser(ctx context.Context, credentials model.Credentials) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 10)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	credentials.Password = string(hashedPassword)

	return s.userRepository.AddUser(ctx, credentials)
}

// FindUserByUsername finds user by username.
func (s userService) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepository.FindUserByUsername(ctx, username)
}

var _ UsersService = (*userService)(nil)

// NewUsersService initiates new instance of UsersService.
func NewUsersService(userRepository repository.UserRepository) UsersService {
	return &userService{
		userRepository: userRepository,
	}
}
