//go:build integration

package repository_integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
	userdomain "github.com/exPriceD/simple-marketplace/internal/domain/user"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/database/postgres"
	userrepo "github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/user"
)

func TestUserRepositoryPG_CRUD(t *testing.T) {
	dsn := os.Getenv("PG_TEST_DSN")
	if dsn == "" {
		t.Skip("PG_TEST_DSN not set")
	}
	ctx := context.Background()
	db, err := postgres.Connect(ctx, postgres.Config{DSN: dsn})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer db.Close()
	if err := postgres.RunMigrations(ctx, db); err != nil {
		t.Fatalf("migrations: %v", err)
	}

	repo := userrepo.NewUserRepositoryPG(db)
	loginVO, _ := userdomain.NewLogin("IntUser_1")
	hashVO, _ := userdomain.NewPasswordHash("argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	u := userdomain.NewUser(shared.RandomIDGenerator{}.NewID(), loginVO, hashVO, time.Now().UTC())

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.FindByLogin(ctx, loginVO)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if got == nil || got.Login().String() != u.Login().String() {
		t.Fatalf("mismatch: %+v", got)
	}
}
