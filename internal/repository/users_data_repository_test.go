package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"e1m0re/passman/internal/models"
	"e1m0re/passman/internal/repository"
)

func Test_usersDataRepository_Add(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type args struct {
		ctx  context.Context
		data models.UsersDataItemInfo
	}
	type want struct {
		usersDataItem *models.UsersDataItem
		err           error
	}
	tests := []struct {
		name string
		mock func() repository.UsersDataRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			args: args{
				ctx: context.Background(),
				data: models.UsersDataItemInfo{
					TypeID:   models.TextItem,
					UserID:   models.UserID(1),
					File:     make([]byte, 0),
					Checksum: make([]byte, 0),
				},
			},
			mock: func() repository.UsersDataRepository {
				repo := repository.NewUsersDataRepository(db)

				mock.
					ExpectQuery("^INSERT INTO users_data_items \\(type, \"user\", file, checksum\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id$").
					WillReturnError(errors.New("something wrong"))

				return repo
			},
			want: want{
				usersDataItem: nil,
				err:           errors.New("something wrong"),
			},
		},
		{
			name: "Successfully case",
			args: args{
				ctx: context.Background(),
				data: models.UsersDataItemInfo{
					TypeID:   models.TextItem,
					UserID:   models.UserID(1),
					File:     make([]byte, 0),
					Checksum: make([]byte, 0),
				},
			},
			mock: func() repository.UsersDataRepository {
				repo := repository.NewUsersDataRepository(db)

				rows := mock.NewRows([]string{"id"}).AddRow(1)

				mock.
					ExpectQuery("^INSERT INTO users_data_items \\(type, \"user\", file, checksum\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id$").
					WillReturnRows(rows)

				return repo
			},
			want: want{
				usersDataItem: &models.UsersDataItem{
					ID:       1,
					TypeID:   models.TextItem,
					UserID:   models.UserID(1),
					File:     make([]byte, 0),
					Checksum: make([]byte, 0),
				},
				err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mock()
			usersDataItem, err := repo.AddItem(test.args.ctx, test.args.data)
			require.Equal(t, test.want.err, err)
			require.Equal(t, test.want.usersDataItem, usersDataItem)
		})
	}
}

func Test_usersDataRepository_FindByID(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type args struct {
		ctx context.Context
		id  models.UsersDataItemID
	}
	type want struct {
		usersDataItem *models.UsersDataItem
		err           error
	}
	tests := []struct {
		name string
		mock func() repository.UsersDataRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() repository.UsersDataRepository {
				repo := repository.NewUsersDataRepository(db)

				mock.
					ExpectQuery("^SELECT \\* FROM users_data_items WHERE id = \\$1 LIMIT 1$").
					WillReturnError(errors.New("something wrong"))

				return repo
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				usersDataItem: nil,
				err:           errors.New("something wrong"),
			},
		},
		{
			name: "UsersDataItem not found",
			mock: func() repository.UsersDataRepository {
				repo := repository.NewUsersDataRepository(db)

				mock.
					ExpectQuery("^SELECT \\* FROM users_data_items WHERE id = \\$1 LIMIT 1$").
					WillReturnError(sql.ErrNoRows)

				return repo
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				usersDataItem: nil,
				err:           repository.ErrorEntityNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() repository.UsersDataRepository {
				repo := repository.NewUsersDataRepository(db)

				rows := sqlxmock.NewRows([]string{"id", "type", "user", "file", "checksum"}).
					AddRow("1", "1", "1", "", "")
				mock.
					ExpectQuery("^SELECT \\* FROM users_data_items WHERE id = \\$1 LIMIT 1$").
					WillReturnRows(rows)

				return repo
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: want{
				usersDataItem: &models.UsersDataItem{
					ID:       1,
					TypeID:   models.TextItem,
					UserID:   1,
					File:     make([]byte, 0),
					Checksum: make([]byte, 0),
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mock()
			usersDataItem, err := repo.FindItemByID(test.args.ctx, test.args.id)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.usersDataItem, usersDataItem)
		})
	}
}

func TestNewUsersDataRepository(t *testing.T) {
	db, _, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}

	repo := repository.NewUsersDataRepository(db)
	assert.Implements(t, (*repository.UsersDataRepository)(nil), repo)
}
