package usecase

import (
	"context"
	"errors"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
	"github.com/exPriceD/simple-marketplace/internal/application/dto/converters"
	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
	"github.com/exPriceD/simple-marketplace/internal/domain/user"
)

type RegisterUserInput struct {
	Login    string
	Password string
}

type RegisterUser struct {
	users  port.UserRepository
	hasher port.PasswordHasher
	clock  port.Clock
	idGen  shared.IDGenerator
}

func NewRegisterUser(users port.UserRepository, hasher port.PasswordHasher, clock port.Clock, idGen shared.IDGenerator) *RegisterUser {
	return &RegisterUser{users: users, hasher: hasher, clock: clock, idGen: idGen}
}

func (uc *RegisterUser) Execute(ctx context.Context, in RegisterUserInput) (dto.UserDTO, error) {
	op := "user.register"

	loginVO, err := user.NewLogin(in.Login)
	if err != nil {
		return dto.UserDTO{}, apperror.New(op, "invalid_login", apperror.KindValidation, err)
	}
	if l := len(in.Password); l < 8 || l > 72 {
		return dto.UserDTO{}, apperror.New(op, "invalid_password", apperror.KindValidation, nil)
	}
	existing, err := uc.users.FindByLogin(ctx, loginVO)
	if err != nil {
		return dto.UserDTO{}, apperror.New(op, "repo_error", apperror.KindInternal, err)
	}
	if existing != nil {
		return dto.UserDTO{}, apperror.New(op, "login_taken", apperror.KindConflict, user.ErrLoginTaken)
	}
	hash, err := uc.hasher.Hash(in.Password)
	if err != nil {
		return dto.UserDTO{}, apperror.New(op, "hash_error", apperror.KindInternal, err)
	}
	hashVO, err := user.NewPasswordHash(hash)
	if err != nil {
		return dto.UserDTO{}, apperror.New(op, "hash_wrap_error", apperror.KindInternal, err)
	}
	u := user.NewUser(uc.idGen.NewID(), loginVO, hashVO, uc.clock.Now())
	if err := uc.users.Create(ctx, u); err != nil {
		if errors.Is(err, user.ErrLoginTaken) {
			return dto.UserDTO{}, apperror.New(op, "login_taken", apperror.KindConflict, err)
		}
		return dto.UserDTO{}, apperror.New(op, "persist_error", apperror.KindInternal, err)
	}
	return converters.UserToDTO(u), nil
}
