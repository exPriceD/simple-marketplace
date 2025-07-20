package userrepo

import (
	"github.com/exPriceD/simple-marketplace/internal/domain/user"
)

// RowToUser преобразует UserRow -> *user.User.
func RowToUser(r UserRow) (*user.User, error) {
	loginVO, err := user.NewLogin(r.Login)
	if err != nil {
		return nil, err
	}
	passHashVO, err := user.NewPasswordHash(r.PasswordHash)
	if err != nil {
		return nil, err
	}
	u := user.RehydrateUser(r.ID, loginVO, passHashVO, r.CreatedAt)
	return u, nil
}

// UserToRow преобразует доменную сущность -> UserRow.
func UserToRow(u *user.User) UserRow {
	return UserRow{
		ID:           u.ID(),
		Login:        u.Login().String(),
		PasswordHash: u.PasswordHash(),
		CreatedAt:    u.CreatedAt(),
	}
}
