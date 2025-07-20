package user

import (
	"time"
)

// User — агрегат пользователя.
type User struct {
	id        string
	login     Login
	passHash  PasswordHash
	createdAt time.Time
}

// NewUser создаёт нового пользователя.
func NewUser(id string, login Login, passHash PasswordHash, createdAt time.Time) *User {
	return &User{
		id:        id,
		login:     login,
		passHash:  passHash,
		createdAt: createdAt.UTC(),
	}
}

// RehydrateUser восстанавливает пользователя из бд.
func RehydrateUser(id string, login Login, passHash PasswordHash, createdAt time.Time) *User {
	return &User{
		id:        id,
		login:     login,
		passHash:  passHash,
		createdAt: createdAt.UTC(),
	}
}

// Геттеры
func (u *User) ID() string           { return u.id }
func (u *User) Login() Login         { return u.login }
func (u *User) PasswordHash() string { return u.passHash.String() }
func (u *User) CreatedAt() time.Time { return u.createdAt }
