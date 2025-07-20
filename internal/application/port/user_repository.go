package port

import (
	"context"

	"github.com/exPriceD/simple-marketplace/internal/domain/user"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	FindByLogin(ctx context.Context, login user.Login) (*user.User, error)
	FindByID(ctx context.Context, id string) (*user.User, error)
}
