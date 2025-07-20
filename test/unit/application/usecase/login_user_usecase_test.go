package usecase_test

import (
	"context"
	"errors"
	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"testing"
	"time"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	userdomain "github.com/exPriceD/simple-marketplace/internal/domain/user"
)

type loginFakeHasher struct{}

func (loginFakeHasher) Hash(p string) (string, error) {
	return "argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", nil
}
func (loginFakeHasher) Verify(h, p string) bool { return p == "Correct123" }

type loginFakeToken struct{}

func (loginFakeToken) Generate(userID, login string, ttl time.Duration) (string, error) {
	return "token123", nil
}
func (loginFakeToken) Parse(token string) (port.Claims, error) {
	return port.Claims{UserID: "uid1", Login: "User_123", Exp: time.Now().Add(time.Hour), Issued: time.Now()}, nil
}

type loginFakeUserRepo struct {
	user *userdomain.User
}

func (r *loginFakeUserRepo) Create(ctx context.Context, u *userdomain.User) error {
	r.user = u
	return nil
}
func (r *loginFakeUserRepo) FindByLogin(ctx context.Context, login userdomain.Login) (*userdomain.User, error) {
	if r.user != nil && r.user.Login().String() == login.String() {
		return r.user, nil
	}
	return nil, nil
}
func (r *loginFakeUserRepo) FindByID(ctx context.Context, id string) (*userdomain.User, error) {
	return nil, nil
}

func TestLoginUser_Success(t *testing.T) {
	loginVO, _ := userdomain.NewLogin("User_123")
	hashVO, _ := userdomain.NewPasswordHash("argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	u := userdomain.NewUser("id1", loginVO, hashVO, time.Now().UTC())

	repo := &loginFakeUserRepo{user: u}
	uc := usecase.NewLoginUser(repo, loginFakeHasher{}, loginFakeToken{})

	out, err := uc.Execute(context.Background(), usecase.LoginUserInput{Login: "User_123", Password: "Correct123"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out.AccessToken != "token123" {
		t.Fatalf("unexpected token")
	}
}

func TestLoginUser_InvalidCreds(t *testing.T) {
	loginVO, _ := userdomain.NewLogin("User_123")
	hashVO, _ := userdomain.NewPasswordHash("argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	u := userdomain.NewUser("id1", loginVO, hashVO, time.Now().UTC())
	repo := &loginFakeUserRepo{user: u}
	uc := usecase.NewLoginUser(repo, loginFakeHasher{}, loginFakeToken{})

	_, err := uc.Execute(context.Background(), usecase.LoginUserInput{Login: "User_123", Password: "Wrong"})
	if err == nil {
		t.Fatalf("expected invalid creds")
	}
	var ae *apperror.AppError
	errors.As(err, &ae)
	if ae.Code != "invalid_credentials" || ae.Kind != apperror.KindAuth {
		t.Fatalf("wrong classification: %+v", ae)
	}
}
