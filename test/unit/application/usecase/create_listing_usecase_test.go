package usecase_test

import (
	"context"
	"errors"
	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"testing"
	"time"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	listingdomain "github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

type fakeListingRepo struct {
	created []*listingdomain.Listing
}

func (r *fakeListingRepo) Create(ctx context.Context, l *listingdomain.Listing) (*listingdomain.Listing, error) {
	cp := *l
	cp = *l.WithID(int64(len(r.created)+1), l.CreatedAt())
	r.created = append(r.created, &cp)
	return &cp, nil
}
func (r *fakeListingRepo) List(ctx context.Context, f port.ListingFilter) ([]*listingdomain.Listing, error) {
	return r.created, nil
}
func (r *fakeListingRepo) Count(ctx context.Context, f port.ListingFilter) (int64, error) {
	return int64(len(r.created)), nil
}

type fakeClock struct{ t time.Time }

func (c fakeClock) Now() time.Time { return c.t }

func TestCreateListing_Success(t *testing.T) {
	repo := &fakeListingRepo{}
	uc := usecase.NewCreateListing(repo, fakeClock{time.Unix(2000, 0).UTC()})
	out, err := uc.Execute(context.Background(), usecase.CreateListingInput{
		AuthorID: "uid1", Title: "Телефон", Description: "Сост 9/10",
		ImageURL: "https://example.com/p.jpg", Price: 1000,
	})
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if out.Title != "Телефон" || out.ID == 0 {
		t.Fatalf("unexpected dto: %+v", out)
	}
}

func TestCreateListing_Validation(t *testing.T) {
	repo := &fakeListingRepo{}
	uc := usecase.NewCreateListing(repo, fakeClock{time.Now().UTC()})
	_, err := uc.Execute(context.Background(), usecase.CreateListingInput{
		AuthorID: "uid1", Title: "", Description: "desc",
		ImageURL: "https://example.com/a.jpg", Price: 10,
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	var ae *apperror.AppError
	errors.As(err, &ae)
	if ae.Kind != apperror.KindValidation {
		t.Fatalf("expected validation kind")
	}
}
