package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/server/repository"
	"github.com/e1m0re/passman/internal/server/service/db/mocks"
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
		mock func() repository.DatumRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			args: args{
				ctx: context.Background(),
				data: model.DatumInfo{
					TypeID:   model.TextItem,
					UserID:   1,
					Metadata: "",
					File:     "",
					Checksum: "",
				},
			},
			mock: func() repository.DatumRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewDatumRepository(mockDBService)

				mock.
					ExpectQuery("^INSERT INTO users_data \\(type, \"user\", metadata, file, checksum\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id$").
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
					UserID:   1,
					Metadata: "",
					File:     "",
					Checksum: "",
				},
			},
			mock: func() repository.DatumRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewDatumRepository(mockDBService)

				rows := mock.NewRows([]string{"id"}).AddRow("1")

				mock.
					ExpectQuery("^INSERT INTO users_data \\(type, \"user\", metadata, file, checksum\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\) RETURNING id$").
					WillReturnRows(rows)

				return repo
			},
			want: want{
				DatumItem: &model.DatumItem{
					ID:       1,
					TypeID:   model.TextItem,
					UserID:   1,
					Metadata: "",
					File:     "",
					Checksum: "",
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
		mock func() repository.DatumRepository
		args args
		want want
	}{
		{
			name: "Something wrong",
			mock: func() repository.DatumRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewDatumRepository(mockDBService)

				mock.
					ExpectQuery("^SELECT \\* FROM users_data WHERE file = \\$1 LIMIT 1$").
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
			mock: func() repository.DatumRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewDatumRepository(mockDBService)

				mock.
					ExpectQuery("^SELECT \\* FROM users_data WHERE file = \\$1 LIMIT 1$").
					WillReturnError(sql.ErrNoRows)

				return repo
			},
			args: args{
				ctx:      context.Background(),
				fileName: "1",
			},
			want: want{
				dataItem: nil,
				err:      repository.ErrorEntityNotFound,
			},
		},
		{
			name: "Successfully case",
			mock: func() repository.DatumRepository {
				mockDBService := mocks.NewDBService(t)
				mockDBService.On("GetDB").Return(db)

				repo := repository.NewDatumRepository(mockDBService)

				rows := sqlxmock.NewRows([]string{"id", "type", "user", "metadata", "file", "checksum"}).
					AddRow("1", "1", "1", "", "1", "")
				mock.
					ExpectQuery("^SELECT \\* FROM users_data WHERE file = \\$1 LIMIT 1$").
					WillReturnRows(rows)

				return repo
			},
			args: args{
				ctx:      context.Background(),
				fileName: "1",
			},
			want: want{
				dataItem: &model.DatumItem{
					ID:       1,
					TypeID:   model.TextItem,
					UserID:   1,
					Metadata: "",
					File:     "1",
					Checksum: "",
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
	mockDBService := mocks.NewDBService(t)

	repo := repository.NewDatumRepository(mockDBService)
	assert.Implements(t, (*repository.DatumRepository)(nil), repo)
}
