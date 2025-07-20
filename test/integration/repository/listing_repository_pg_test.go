//go:build integration

package repository_integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
	listingdomain "github.com/exPriceD/simple-marketplace/internal/domain/listing"
	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
	userdomain "github.com/exPriceD/simple-marketplace/internal/domain/user"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/database/postgres"
	listingrepo "github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/listing"
	userrepo "github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/user"
)

func TestListingRepositoryPG_ListFilter(t *testing.T) {
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
		t.Fatalf("migrate: %v", err)
	}

	userRepo := userrepo.NewUserRepositoryPG(db)
	loginVO, _ := userdomain.NewLogin("Lister_1")
	hashVO, _ := userdomain.NewPasswordHash("argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	u := userdomain.NewUser(shared.RandomIDGenerator{}.NewID(), loginVO, hashVO, time.Now().UTC())
	if err := userRepo.Create(ctx, u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	repo := listingrepo.NewListingRepositoryPG(db)

	add := func(title string, price int64) {
		titleVO, _ := listingdomain.NewTitle(title)
		descVO, _ := listingdomain.NewDescription("Desc")
		imgVO, _ := listingdomain.NewImageURL("https://example.com/" + title + ".png")
		priceVO, _ := listingdomain.NewPrice(price)
		l := listingdomain.NewListing(titleVO, descVO, imgVO, priceVO, u.ID(), time.Now().UTC())
		if _, err := repo.Create(ctx, l); err != nil {
			t.Fatalf("create listing: %v", err)
		}
	}
	add("A", 100)
	add("B", 200)
	add("C", 300)

	minPrice := int64(150)
	maxPrice := int64(250)
	items, err := repo.List(ctx, port.ListingFilter{
		PriceMin: &minPrice,
		PriceMax: &maxPrice,
		SortBy:   "price",
		SortDir:  "asc",
		Limit:    10,
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 || items[0].Price() != 200 {
		t.Fatalf("unexpected filtered result %+v", items)
	}
}
