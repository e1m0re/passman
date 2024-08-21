package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"e1m0re/passman/internal/models"
)

type usersDataRepository struct {
	db *sqlx.DB
}

// Add creates new users data item.
func (repo usersDataRepository) Add(ctx context.Context, data models.UsersDataItemInfo) (*models.UsersDataItem, error) {
	usersDataItem := &models.UsersDataItem{
		TypeID:   data.TypeID,
		UserID:   data.UserID,
		File:     data.File,
		Checksum: data.Checksum,
	}

	query := "INSERT INTO users_data_items (type, \"user\", file, checksum) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRowxContext(ctx, query, data.TypeID, data.UserID, data.File, data.Checksum).Scan(&usersDataItem.ID)
	if err != nil {
		return nil, err
	}

	return usersDataItem, nil
}

// FindByID finds and returns users data item.
func (repo usersDataRepository) FindByID(ctx context.Context, id models.UsersDataItemID) (*models.UsersDataItem, error) {
	usersDataItem := &models.UsersDataItem{}
	query := "SELECT * FROM users_data_items WHERE id = $1 LIMIT 1"
	err := repo.db.GetContext(ctx, usersDataItem, query, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrorEntityNotFound
	case err != nil:
		return nil, err
	default:
		return usersDataItem, nil
	}
}

var _ UsersDataRepository = (*usersDataRepository)(nil)

// NewUsersDataRepository initiates new instance of UsersDataRepository.
func NewUsersDataRepository(db *sqlx.DB) UsersDataRepository {
	return &usersDataRepository{
		db: db,
	}
}
