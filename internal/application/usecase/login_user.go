package usecase

import (
	"context"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
	"github.com/exPriceD/simple-marketplace/internal/application/dto/converters"
	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/domain/user"
)

type LoginUserInput struct {
	Login    string
	Password string
}

type LoginUser struct {
	users  port.UserRepository
	hasher port.PasswordHasher
	tokens port.TokenService
}

func NewLoginUser(users port.UserRepository, hasher port.PasswordHasher, tokens port.TokenService) *LoginUser {
	return &LoginUser{users: users, hasher: hasher, tokens: tokens}
}

func (uc *LoginUser) Execute(ctx context.Context, in LoginUserInput) (dto.AuthDTO, error) {
	op := "user.login"
	loginVO, err := user.NewLogin(in.Login)
	if err != nil {
		return dto.AuthDTO{}, apperror.New(op, "invalid_credentials", apperror.KindAuth, err)
	}
	u, err := uc.users.FindByLogin(ctx, loginVO)
	if err != nil {
		return dto.AuthDTO{}, apperror.New(op, "repo_error", apperror.KindInternal, err)
	}
	if u == nil {
		return dto.AuthDTO{}, apperror.New(op, "invalid_credentials", apperror.KindAuth, nil)
	}
	if !uc.hasher.Verify(u.PasswordHash(), in.Password) {
		return dto.AuthDTO{}, apperror.New(op, "invalid_credentials", apperror.KindAuth, nil)
	}
	token, err := uc.tokens.Generate(u.ID(), u.Login().String(), 0)
	if err != nil {
		return dto.AuthDTO{}, apperror.New(op, "token_issue_error", apperror.KindInternal, err)
	}
	return dto.AuthDTO{
		AccessToken: token,
		User:        converters.UserToDTO(u),
	}, nil
}
