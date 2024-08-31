package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/e1m0re/passman/internal/model"
)

func Test_datumRepository_AddItem(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type args struct {
		ctx  context.Context
		data model.DatumInfo
	}
	type want struct {
		DatumItem *model.DatumItem
		err       error
	}
	tests := []struct {
		name string
		mock func() DatumRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			args: args{
				ctx: context.Background(),
				data: model.DatumInfo{
					TypeID:   model.TextItem,
					UserID:   model.UserID(1),
					File:     make([]byte, 0),
					Checksum: make([]byte, 0),
				},
			},
			mock: func() DatumRepository {
				repo := NewDatumRepository(db)

				mock.
					ExpectQuery("^INSERT INTO users_data \\(type, user, file, checksum\\) VALUES \\(\\?,\\?,\\?,\\?\\) RETURNING id$").
					WillReturnError(errors.New("something wrong"))

				return repo
			},
			want: want{
				DatumItem: nil,
				err:       errors.New("something wrong"),
			},
		},
		{
			name: "Successfully case",
			args: args{
				ctx: context.Background(),
				data: model.DatumInfo{
					TypeID:   model.TextItem,
					UserID:   model.UserID(1),
					File:     make([]byte, 0),
					Checksum: make([]byte, 0),
				},
			},
			mock: func() DatumRepository {
				repo := NewDatumRepository(db)

				rows := mock.NewRows([]string{"id"}).AddRow(model.DatumID("1"))

				mock.
					ExpectQuery("^INSERT INTO users_data \\(type, user, file, checksum\\) VALUES \\(\\?,\\?,\\?,\\?\\) RETURNING id$").
					WillReturnRows(rows)

				return repo
			},
			want: want{
				DatumItem: &model.DatumItem{
					ID:       model.DatumID("1"),
					TypeID:   model.TextItem,
					UserID:   model.UserID(1),
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
			require.Equal(t, test.want.DatumItem, usersDataItem)
		})
	}
}

func Test_datumRepository_FindItemByFileName(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	type args struct {
		ctx      context.Context
		fileName string
	}
	type want struct {
		dataItem *model.DatumItem
		err      error
	}
	tests := []struct {
		name string
		mock func() DatumRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() DatumRepository {
				repo := NewDatumRepository(db)

				mock.
					ExpectQuery("^SELECT \\* FROM users_data WHERE file = \\? LIMIT 1$").
					WillReturnError(errors.New("something wrong"))

				return repo
			},
			args: args{
				ctx:      context.Background(),
				fileName: "1",
			},
			want: want{
				dataItem: nil,
				err:      errors.New("something wrong"),
			},
		},
		{
			name: "UsersDataItem not found",
			mock: func() DatumRepository {
				repo := NewDatumRepository(db)

				mock.
					ExpectQuery("^SELECT \\* FROM users_data WHERE file = \\? LIMIT 1$").
					WillReturnError(sql.ErrNoRows)

				return repo
			},
			args: args{
				ctx:      context.Background(),
				fileName: "1",
			},
			want: want{
				dataItem: nil,
				err:      ErrorEntityNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() DatumRepository {
				repo := NewDatumRepository(db)

				rows := sqlxmock.NewRows([]string{"id", "type", "user", "file", "checksum"}).
					AddRow("1", "1", "1", "1", "")
				mock.
					ExpectQuery("^SELECT \\* FROM users_data WHERE file = \\? LIMIT 1$").
					WillReturnRows(rows)

				return repo
			},
			args: args{
				ctx:      context.Background(),
				fileName: "1",
			},
			want: want{
				dataItem: &model.DatumItem{
					ID:       model.DatumID("1"),
					TypeID:   model.TextItem,
					UserID:   1,
					File:     []byte("1"),
					Checksum: make([]byte, 0),
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := test.mock()
			usersDataItem, err := repo.FindItemByFileName(test.args.ctx, test.args.fileName)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.dataItem, usersDataItem)
		})
	}
}

func Test_datumRepository_NewDatumRepository(t *testing.T) {
	db, _, err := sqlxmock.Newx()
	if err != nil {
		panic(err)
	}

	repo := NewDatumRepository(db)
	assert.Implements(t, (*DatumRepository)(nil), repo)
}
