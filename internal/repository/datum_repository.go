package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/service/db"
)

type datumRepository struct {
	db db.DBService
}

// AddItem creates new users data item.
func (repo datumRepository) AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error) {
	datumItem := &model.DatumItem{
		UserID:   datumInfo.UserID,
		TypeID:   datumInfo.TypeID,
		File:     datumInfo.File,
		Checksum: datumInfo.Checksum,
	}

	query := "INSERT INTO users_data (type, \"user\", file, checksum) VALUES ($1,$2,$3,$4) RETURNING id"
	err := repo.db.GetDB().
		QueryRowxContext(ctx, query, datumInfo.TypeID, datumInfo.UserID, datumInfo.File, datumInfo.Checksum).
		Scan(&datumItem.ID)
	if err != nil {
		return nil, err
	}

	return datumItem, nil
}

// FindItemByFileName finds and returns users data item by filename.
func (repo datumRepository) FindItemByFileName(ctx context.Context, fileName string) (*model.DatumItem, error) {
	datumItem := &model.DatumItem{}
	query := "SELECT * FROM users_data WHERE file = $1 LIMIT 1"
	err := repo.db.GetDB().GetContext(ctx, datumItem, query, fileName)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrorEntityNotFound
	case err != nil:
		return nil, err
	default:
		return datumItem, nil
	}
}

var _ DatumRepository = (*datumRepository)(nil)

// NewDatumRepository initiates new instance of DatumRepository.
func NewDatumRepository(db db.DBService) DatumRepository {
	return &datumRepository{db: db}
}
