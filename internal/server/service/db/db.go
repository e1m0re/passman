package db

import (
	"io"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=DBService
type DBService interface {
	io.Closer
	// GetDB returns link to db instance.
	GetDB() *sqlx.DB
}

type dbService struct {
	*sqlx.DB
}

// Close closes the connection to the dbService.
func (d dbService) Close() error {
	return d.DB.Close()
}

// GetDB returns link to db instance.
func (d dbService) GetDB() *sqlx.DB {
	return d.DB
}

var _ DBService = (*dbService)(nil)

// NewDBService initiates new instance of DBService.
func NewDBService(conf DatabaseConfig) (DBService, error) {
	db, err := sqlx.Open(conf.Driver, conf.Url)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * time.Duration(conf.ConnMaxLifetimeInMinute))
	db.SetMaxOpenConns(conf.MaxOpenConnections)
	db.SetMaxIdleConns(conf.MaxIdleConnections)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &dbService{
		db,
	}, nil
}
