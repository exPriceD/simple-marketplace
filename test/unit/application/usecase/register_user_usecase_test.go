package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	userdomain "github.com/exPriceD/simple-marketplace/internal/domain/user"
)

type fakeUserRepo struct {
	usersByLogin map[string]*userdomain.User
	errFind      error
	errCreate    error
}

func (f *fakeUserRepo) Create(ctx context.Context, u *userdomain.User) error {
	if f.errCreate != nil {
		return f.errCreate
	}
	if _, exists := f.usersByLogin[u.Login().String()]; exists {
		return userdomain.ErrLoginTaken
	}
	f.usersByLogin[u.Login().String()] = u
	return nil
}
func (f *fakeUserRepo) FindByLogin(ctx context.Context, login userdomain.Login) (*userdomain.User, error) {
	if f.errFind != nil {
		return nil, f.errFind
	}
	return f.usersByLogin[login.String()], nil
}
func (f *fakeUserRepo) FindByID(ctx context.Context, id string) (*userdomain.User, error) {
	return nil, nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(p string) (string, error) {
	return "argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", nil
}
func (fakeHasher) Verify(h, p string) bool { return true }

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

type seqIDGen struct{ n int }

func (g *seqIDGen) NewID() string {
	g.n++
	return "id" + strconvI(g.n)
}
func strconvI(i int) string {
	return fmt.Sprintf("%d", i)
}

func TestRegisterUser_Success(t *testing.T) {
	repo := &fakeUserRepo{usersByLogin: make(map[string]*userdomain.User)}
	uc := usecase.NewRegisterUser(repo, fakeHasher{}, fixedClock{t: time.Unix(1000, 0).UTC()}, &seqIDGen{})

	out, err := uc.Execute(context.Background(), usecase.RegisterUserInput{
		Login: "User_123", Password: "StrongPass1",
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out.Login != "User_123" || out.ID == "" {
		t.Fatalf("unexpected dto: %+v", out)
	}
}

func TestRegisterUser_LoginTaken(t *testing.T) {
	repo := &fakeUserRepo{usersByLogin: make(map[string]*userdomain.User)}
	uc := usecase.NewRegisterUser(repo, fakeHasher{}, fixedClock{time.Now().UTC()}, &seqIDGen{})
	_, _ = uc.Execute(context.Background(), usecase.RegisterUserInput{Login: "dupUser", Password: "StrongPass1"})
	_, err := uc.Execute(context.Background(), usecase.RegisterUserInput{Login: "dupUser", Password: "StrongPass1"})
	if err == nil {
		t.Fatalf("expected error")
	}
	var ae *apperror.AppError
	ok := errors.As(err, &ae)
	if !ok || ae.Code != "login_taken" || ae.Kind != apperror.KindConflict {
		t.Fatalf("wrong app error: %+v", err)
	}
}

func TestRegisterUser_InvalidLogin(t *testing.T) {
	repo := &fakeUserRepo{usersByLogin: make(map[string]*userdomain.User)}
	uc := usecase.NewRegisterUser(repo, fakeHasher{}, fixedClock{time.Now().UTC()}, &seqIDGen{})
	_, err := uc.Execute(context.Background(), usecase.RegisterUserInput{Login: "x", Password: "StrongPass1"})
	if err == nil {
		t.Fatalf("expected invalid login error")
	}
	var ae *apperror.AppError
	errors.As(err, &ae)
	if ae.Code != "invalid_login" || ae.Kind != apperror.KindValidation {
		t.Fatalf("wrong error classification: %+v", ae)
	}
}

func TestRegisterUser_InvalidPassword(t *testing.T) {
	repo := &fakeUserRepo{usersByLogin: make(map[string]*userdomain.User)}
	uc := usecase.NewRegisterUser(repo, fakeHasher{}, fixedClock{time.Now().UTC()}, &seqIDGen{})
	_, err := uc.Execute(context.Background(), usecase.RegisterUserInput{Login: "User_123", Password: "short"})
	if err == nil {
		t.Fatalf("expected invalid password")
	}
	var ae *apperror.AppError
	errors.As(err, &ae)
	if ae.Code != "invalid_password" {
		t.Fatalf("wrong code: %s", ae.Code)
	}
}
