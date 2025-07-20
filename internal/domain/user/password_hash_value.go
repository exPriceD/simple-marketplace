package user

import "strings"

type PasswordHash struct {
	value string
}

func NewPasswordHash(raw string) (PasswordHash, error) {
	if !strings.HasPrefix(raw, "argon2") {
		return PasswordHash{}, ErrInvalidPasswordHash
	}
	return PasswordHash{value: raw}, nil
}

var ErrInvalidPasswordHash = ErrInvalidPasswordHashType{}

type ErrInvalidPasswordHashType struct{}

func (ErrInvalidPasswordHashType) Error() string { return "invalid_password_hash" }

func (h PasswordHash) String() string { return h.value }
