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

// Add creates new user.
func (repo userRepository) Add(ctx context.Context, credentials models.Credentials) (*models.User, error) {
	user := &models.User{
		Username: credentials.Username,
		Password: credentials.Password,
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRowxContext(ctx, query, credentials.Username, credentials.Password).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && err.(*pgconn.PgError).Code == "23505" {
			return nil, ErrorBusyLogin
		}
		return nil, err
	}

	return user, nil
}

// FindByID finds and returns user instance by id.
func (repo userRepository) FindByID(ctx context.Context, id models.UserID) (*models.User, error) {
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

// FindByUsername finds and returns user instance by username.
func (repo userRepository) FindByUsername(ctx context.Context, username []byte) (*models.User, error) {
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
