package userrepo

import "time"

// UserRow — структура уровня хранилища.
type UserRow struct {
	ID        string    `db:"id"`
	Login     string    `db:"login"`
	PassHash  string    `db:"pass_hash"`
	CreatedAt time.Time `db:"created_at"`
}
