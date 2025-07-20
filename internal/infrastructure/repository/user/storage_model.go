package userrepo

import "time"

// UserRow — структура уровня хранилища.
type UserRow struct {
	ID           string    `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
