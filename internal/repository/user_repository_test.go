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

	"e1m0re/passman/internal/models"
	"e1m0re/passman/internal/repository"
)

var (
	errorSomethingWrong = errors.New("something wrong")
)

func Test_userRepository_AddUser(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type args struct {
		ctx      context.Context
		userInfo models.UserInfo
	}
	type want struct {
		user *models.User
		err  error
	}
	tests := []struct {
		name string
		mock func() repository.UserRepository
		args args
		want want
	}{
		{
			name: "Login is busy",
			args: args{
				ctx: context.Background(),
				userInfo: models.UserInfo{
					Username: "username",
					Password: "password",
				},
			},
			want: want{
				user: nil,
				err:  errorSomethingWrong,
			},
			mock: func() repository.UserRepository {
				repo := repository.NewUserRepository(db)

				mock.
					ExpectQuery("^INSERT INTO users \\(username, password\\) VALUES \\(\\$1, \\$2\\) RETURNING id$").
					WillReturnError(errorSomethingWrong)

				return repo
			},
		},
		{
			name: "Login is busy",
			args: args{
				ctx: context.Background(),
				userInfo: models.UserInfo{
					Username: "username",
					Password: "password",
				},
			},
			want: want{
				user: nil,
				err:  repository.ErrorBusyLogin,
			},
			mock: func() repository.UserRepository {
				repo := repository.NewUserRepository(db)

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
				userInfo: models.UserInfo{
					Username: "username",
					Password: "password",
				},
			},
			want: want{
				user: &models.User{
					ID:       1,
					Username: "username",
					Password: "password",
				},
				err: nil,
			},
			mock: func() repository.UserRepository {
				repo := repository.NewUserRepository(db)

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
			user, err := repo.AddUser(test.args.ctx, test.args.userInfo)
			require.Equal(t, test.want.err, err)
			require.Equal(t, test.want.user, user)
		})
	}
}

func Test_userRepository_FindUserById(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
		mock func() repository.UserRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() repository.UserRepository {
				repo := repository.NewUserRepository(db)

				mock.
					ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1 LIMIT 1$").
					WillReturnError(errorSomethingWrong)

				return repo
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				user: nil,
				err:  errorSomethingWrong,
			},
		},
		{
			name: "User not found",
			mock: func() repository.UserRepository {
				repo := repository.NewUserRepository(db)

				mock.
					ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1 LIMIT 1$").
					WillReturnError(sql.ErrNoRows)

				return repo
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				user: nil,
				err:  repository.ErrorEntityNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() repository.UserRepository {
				repo := repository.NewUserRepository(db)

				rows := sqlxmock.NewRows([]string{"id", "username", "password"}).
					AddRow("1", "username", "password")
				mock.
					ExpectQuery("^SELECT \\* FROM users WHERE id = \\$1 LIMIT 1$").
					WillReturnRows(rows)

				return repo
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
			repo := test.mock()
			user, err := repo.FindUserById(test.args.ctx, test.args.id)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.user, user)
		})
	}
}

func Test_userRepository_FindUserByUsername(t *testing.T) {
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
		user *models.User
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
				repo := repository.NewUserRepository(db)

				mock.
					ExpectQuery("^SELECT \\* FROM users WHERE username = \\$1 LIMIT 1$").
					WillReturnError(errorSomethingWrong)

				return repo
			},
			args: args{
				ctx:      context.Background(),
				username: "username",
			},
			want: want{
				user: nil,
				err:  errorSomethingWrong,
			},
		},
		{
			name: "User not found",
			mock: func() repository.UserRepository {
				repo := repository.NewUserRepository(db)

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
				repo := repository.NewUserRepository(db)

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
			repo := test.mock()
			user, err := repo.FindUserByUsername(test.args.ctx, test.args.username)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.user, user)
		})
	}
}

func TestNewUserRepository(t *testing.T) {
	db, _, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}

	repo := repository.NewUserRepository(db)
	assert.Implements(t, (*repository.UserRepository)(nil), repo)
}
