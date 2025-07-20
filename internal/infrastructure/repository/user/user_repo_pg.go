package userrepo

import (
	"context"
	"errors"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/domain/user"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepositoryPG — реализация порта UserRepository для Postgres.
type UserRepositoryPG struct {
	db *postgres.DB
}

func NewUserRepositoryPG(db *postgres.DB) *UserRepositoryPG {
	return &UserRepositoryPG{db: db}
}

func (r *UserRepositoryPG) Create(ctx context.Context, u *user.User) error {
	row := UserToRow(u)
	_, err := r.db.Pool.Exec(ctx, `
		INSERT INTO users (id, login, pass_hash, created_at)
		VALUES ($1,$2,$3,$4)
	`, row.ID, row.Login, row.PassHash, row.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return user.ErrLoginTaken
		}
		return err
	}
	return nil
}

func (r *UserRepositoryPG) FindByLogin(ctx context.Context, login user.Login) (*user.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, login, pass_hash, created_at
		FROM users
		WHERE login=$1
	`, login.String())

	var ur UserRow
	if err := row.Scan(&ur.ID, &ur.Login, &ur.PassHash, &ur.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return RowToUser(ur)
}

func (r *UserRepositoryPG) FindByID(ctx context.Context, id string) (*user.User, error) {
	row := r.db.Pool.QueryRow(ctx, `
		SELECT id, login, pass_hash, created_at
		FROM users
		WHERE id=$1
	`, id)

	var ur UserRow
	if err := row.Scan(&ur.ID, &ur.Login, &ur.PassHash, &ur.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return RowToUser(ur)
}

var _ port.UserRepository = (*UserRepositoryPG)(nil)

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
