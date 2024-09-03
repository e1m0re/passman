package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/server/service/db"
)

var (
	ErrorBusyLogin      = errors.New("busy login")
	ErrorEntityNotFound = errors.New("entity not found")
)

type userRepository struct {
	db db.DBService
}

// AddUser creates new user.
func (repo userRepository) AddUser(ctx context.Context, credentials model.Credentials) (*model.User, error) {
	user := &model.User{
		Username: credentials.Username,
		Password: credentials.Password,
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err := repo.db.GetDB().QueryRowxContext(ctx, query, credentials.Username, credentials.Password).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && err.(*pgconn.PgError).Code == "23505" {
			return nil, ErrorBusyLogin
		}
		return nil, err
	}

	return user, nil
}

// FindUserByUsername finds and returns user instance by username.
func (repo userRepository) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	err := repo.db.GetDB().GetContext(ctx, user, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrorEntityNotFound
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

var _ UserRepository = (*userRepository)(nil)

// NewUserRepository initiates new instance of UserRepository.
func NewUserRepository(db db.DBService) UserRepository {
	return &userRepository{
		db: db,
	}
}
