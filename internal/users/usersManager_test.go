package users_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	errors2 "e1m0re/passman/internal/errors"
	"e1m0re/passman/internal/models"
	"e1m0re/passman/internal/repository"
	"e1m0re/passman/internal/repository/mocks"
	"e1m0re/passman/internal/users"
)

func Test_usersManager_FindUserByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  models.UserID
	}
	type want struct {
		user *models.User
		err  error
	}
	tests := []struct {
		name string
		mock func() users.Manager
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() users.Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, errors.New("something wrong"))

				return users.NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				user: nil,
				err:  errors.New("something wrong"),
			},
		},
		{
			name: "User not found",
			mock: func() users.Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, repository.ErrorEntityNotFound)

				return users.NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				user: nil,
				err:  errors2.ErrorUserNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() users.Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.User{
						ID:       1,
						Username: "username",
						Password: "password",
					}, nil)

				return users.NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				user: &models.User{
					ID:       1,
					Username: "username",
					Password: "password",
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mgr := test.mock()
			user, err := mgr.FindUserByID(test.args.ctx, test.args.id)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.user, user)
		})
	}
}

func Test_usersManager_FindUserByUsername(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}
	type want struct {
		user *models.User
		err  error
	}
	tests := []struct {
		name string
		mock func() users.Manager
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() users.Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("something wrong"))

				return users.NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			want: want{
				user: nil,
				err:  errors.New("something wrong"),
			},
		},
		{
			name: "User not found",
			mock: func() users.Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(nil, repository.ErrorEntityNotFound)

				return users.NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			want: want{
				user: nil,
				err:  errors2.ErrorUserNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() users.Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByUsername", mock.Anything, mock.AnythingOfType("string")).
					Return(&models.User{
						ID:       1,
						Username: "username",
						Password: "password",
					}, nil)

				return users.NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			want: want{
				user: &models.User{
					ID:       1,
					Username: "username",
					Password: "password",
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mgr := test.mock()
			user, err := mgr.FindUserByUsername(test.args.ctx, test.args.username)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.user, user)
		})
	}
}

func TestNewUsersManager(t *testing.T) {
	mockUsersRepository := mocks.NewUserRepository(t)

	mgr := users.NewUsersManager(mockUsersRepository)
	assert.Implements(t, (*users.Manager)(nil), mgr)
}
