package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/server/repository"
	"github.com/e1m0re/passman/internal/server/service/db/mocks"
)

func Test_userRepository_Add(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type args struct {
		ctx         context.Context
		credentials model.Credentials
	}
	type want struct {
		user *model.User
		err  error
	}
	tests := []struct {
		name string
		mock func() repository.UserRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			args: args{
				ctx: context.Background(),
				credentials: model.Credentials{
					Username: "username",
					Password: "password",
				},
			},
			want: want{
				user: nil,
				err:  errors.New("something wrong"),
			},
			mock: func() repository.UserRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewUserRepository(mockDBService)

				mock.
					ExpectQuery("^INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id$").
					WillReturnError(errors.New("something wrong"))

				return repo
			},
		},
		{
			name: "Login is busy",
			args: args{
				ctx: context.Background(),
				credentials: model.Credentials{
					Username: "username",
					Password: "password",
				},
			},
			want: want{
				user: nil,
				err:  repository.ErrorBusyLogin,
			},
			mock: func() repository.UserRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewUserRepository(mockDBService)

				mock.
					ExpectQuery("^INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id$").
					WillReturnError(&pgconn.PgError{Code: "23505"})

				return repo
			},
		},
		{
			name: "Successfully case",
			args: args{
				ctx: context.Background(),
				credentials: model.Credentials{
					Username: "username",
					Password: "password",
				},
			},
			want: want{
				user: &model.User{
					ID:       1,
					Username: "username",
					Password: "password",
				},
				err: nil,
			},
			mock: func() repository.UserRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewUserRepository(mockDBService)

				rows := mock.NewRows([]string{"id"}).AddRow(1)

				mock.
					ExpectQuery("^INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id$").
					WillReturnRows(rows)

				return repo
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mock()
			user, err := repo.AddUser(test.args.ctx, test.args.credentials)
			require.Equal(t, test.want.err, err)
			require.Equal(t, test.want.user, user)
		})
	}
}

func Test_userRepository_FindByUsername(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type args struct {
		ctx      context.Context
		username string
	}
	type want struct {
		user *model.User
		err  error
	}
	tests := []struct {
		name string
		mock func() repository.UserRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() repository.UserRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewUserRepository(mockDBService)

				mock.
					ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1 LIMIT 1$").
					WillReturnError(errors.New("something wrong"))

				return repo
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
			mock: func() repository.UserRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewUserRepository(mockDBService)

				mock.
					ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1 LIMIT 1$").
					WillReturnError(sql.ErrNoRows)

				return repo
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			want: want{
				user: nil,
				err:  repository.ErrorEntityNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() repository.UserRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewUserRepository(mockDBService)

				rows := sqlxmock.NewRows([]string{"id", "username", "password"}).
					AddRow("1", "username", "password")
				mock.
					ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1 LIMIT 1$").
					WillReturnRows(rows)

				return repo
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			want: want{
				user: &model.User{
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
			repo := test.mock()
			user, err := repo.FindUserByUsername(test.args.ctx, test.args.username)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.user, user)
		})
	}
}

func TestNewUserRepository(t *testing.T) {
	mockDBService := mocks.NewDBService(t)

	repo := repository.NewUserRepository(mockDBService)
	assert.Implements(t, (*repository.UserRepository)(nil), repo)
}
