package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"e1m0re/passman/internal/models"
)

var (
	ErrorBusyLogin      = errors.New("busy login")
	ErrorEntityNotFound = errors.New("entity not found")
)

type userRepository struct {
	db *sqlx.DB
}

// AddUser creates new user.
func (repo userRepository) AddUser(ctx context.Context, userInfo models.UserInfo) (*models.User, error) {
	user := &models.User{
		Username: userInfo.Username,
		Password: userInfo.Password,
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRowxContext(ctx, query, userInfo.Username, userInfo.Password).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && err.(*pgconn.PgError).Code == "23505" {
			return nil, ErrorBusyLogin
		}
		return nil, err
	}

	return user, nil
}

// FindUserById finds and returns user instance by id or nil.
func (repo userRepository) FindUserById(ctx context.Context, id models.UserID) (*models.User, error) {
	user := &models.User{}
	query := "SELECT * FROM users WHERE id = $1 LIMIT 1"
	err := repo.db.GetContext(ctx, user, query, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrorEntityNotFound
	case err != nil:
		return nil, err
	default:
		return user, nil
	}
}

// FindUserByUsername finds and returns user instance by username or nil.
func (repo userRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT * FROM users WHERE username = $1 LIMIT 1"
	err := repo.db.GetContext(ctx, user, query, username)
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
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
