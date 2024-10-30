package postgresrepo

import (
	"context"
	"errors"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

// AddUser add user to repository.
// Returns uuid of added user or error if user creation failed.
// Login must be unique, otherwise ErrUserAlreadyExists will be returned.
func (p PostgresRepository) AddUser(
	ctx context.Context,
	login, passHash string,
) (uuid.UUID, error) {
	row := p.pool.QueryRow(
		ctx,
		queryAddUser,
		login,
		passHash,
	)

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, apperrors.ErrUserAlreadyExists
		}
		return uuid.Nil, err
	}

	return id, nil
}

// GetUserByLogin retrieves a user from the PostgreSQL database using the given login.
// It executes a query to fetch the user's ID, login, and password hash.
// Returns a models.User if the query is successful.
// Returns ErrWrongCredentials if no user is found with the specified login.
// Returns an error for any other query failure.
func (p PostgresRepository) GetUserByLogin(
	ctx context.Context,
	login string,
) (models.User, error) {
	row := p.pool.QueryRow(
		ctx,
		queryGetUser,
		login,
	)

	var user models.User
	err := row.Scan(&user.ID, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, apperrors.ErrWrongCredentials
		}
		return models.User{}, err
	}

	return user, nil
}
