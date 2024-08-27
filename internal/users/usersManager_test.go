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
)

func Test_usersManager_AddUser(t *testing.T) {
	type args struct {
		ctx        context.Context
		credential models.Credentials
	}
	type want struct {
		user *models.User
		err  error
	}
	tests := []struct {
		name string
		mock func() Manager
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("AddUser", mock.Anything, mock.AnythingOfType("models.Credentials")).
					Return(nil, errors.New("something wrong"))

				return NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx: context.Background(),
				credential: models.Credentials{
					Password: []byte("username"),
					Username: []byte("password"),
				},
			},
			want: want{
				user: nil,
				err:  errors.New("something wrong"),
			},
		},
		{
			name: "User not found",
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("AddUser", mock.Anything, mock.AnythingOfType("models.Credentials")).
					Return(nil, repository.ErrorBusyLogin)

				return NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx: context.Background(),
				credential: models.Credentials{
					Password: []byte("username"),
					Username: []byte("password"),
				},
			},
			want: want{
				user: nil,
				err:  repository.ErrorBusyLogin,
			},
		},
		{
			name: "Successfully case",
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("AddUser", mock.Anything, mock.AnythingOfType("models.Credentials")).
					Return(&models.User{
						ID:       1,
						Username: []byte("username"),
						Password: []byte("password"),
					}, nil)

				return NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx: context.Background(),
				credential: models.Credentials{
					Password: []byte("username"),
					Username: []byte("password"),
				},
			},
			want: want{
				user: &models.User{
					ID:       1,
					Username: []byte("username"),
					Password: []byte("password"),
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mgr := test.mock()
			user, err := mgr.AddUser(test.args.ctx, test.args.credential)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.user, user)
		})
	}
}

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
		mock func() Manager
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, errors.New("something wrong"))

				return NewUsersManager(mockUsersRepository)
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
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(nil, repository.ErrorEntityNotFound)

				return NewUsersManager(mockUsersRepository)
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
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByID", mock.Anything, mock.AnythingOfType("models.UserID")).
					Return(&models.User{
						ID:       1,
						Username: []byte("username"),
						Password: []byte("password"),
					}, nil)

				return NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				user: &models.User{
					ID:       1,
					Username: []byte("username"),
					Password: []byte("password"),
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
		username []byte
	}
	type want struct {
		user *models.User
		err  error
	}
	tests := []struct {
		name string
		mock func() Manager
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByUsername", mock.Anything, []byte("username")).
					Return(nil, errors.New("something wrong"))

				return NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx:      context.Background(),
				username: []byte("username"),
			},
			want: want{
				user: nil,
				err:  errors.New("something wrong"),
			},
		},
		{
			name: "User not found",
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByUsername", mock.Anything, []byte("username")).
					Return(nil, repository.ErrorEntityNotFound)

				return NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx:      context.Background(),
				username: []byte("username"),
			},
			want: want{
				user: nil,
				err:  errors2.ErrorUserNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() Manager {
				mockUsersRepository := mocks.NewUserRepository(t)
				mockUsersRepository.
					On("FindUserByUsername", mock.Anything, []byte("username")).
					Return(&models.User{
						ID:       1,
						Username: []byte("username"),
						Password: []byte("password"),
					}, nil)

				return NewUsersManager(mockUsersRepository)
			},
			args: args{
				ctx:      context.Background(),
				username: []byte("username"),
			},
			want: want{
				user: &models.User{
					ID:       1,
					Username: []byte("username"),
					Password: []byte("password"),
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

	mgr := NewUsersManager(mockUsersRepository)
	assert.Implements(t, (*Manager)(nil), mgr)
}
