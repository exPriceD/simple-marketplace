package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/application/usecase"
	listingdomain "github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

type fakeListRepo struct {
	items []*listingdomain.Listing
}

func (r *fakeListRepo) Create(ctx context.Context, l *listingdomain.Listing) (*listingdomain.Listing, error) {
	return nil, nil
}
func (r *fakeListRepo) List(ctx context.Context, f port.ListingFilter) ([]*listingdomain.Listing, error) {
	return r.items, nil
}
func (r *fakeListRepo) Count(ctx context.Context, f port.ListingFilter) (int64, error) {
	return int64(len(r.items)), nil
}

func TestListListings_Success(t *testing.T) {
	repo := &fakeListRepo{}
	titleVO, _ := listingdomain.NewTitle("Item1")
	descVO, _ := listingdomain.NewDescription("Desc")
	imgVO, _ := listingdomain.NewImageURL("https://example.com/i.png")
	priceVO, _ := listingdomain.NewPrice(500)
	created := time.Now().UTC()
	orig := listingdomain.RehydrateListing(42, titleVO, descVO, imgVO, priceVO, "user123", "user123", created)
	repo.items = []*listingdomain.Listing{
		orig,
	}

	uc := usecase.NewListListings(repo)
	out, err := uc.Execute(context.Background(), usecase.ListListingsInput{
		Limit: 10, Offset: 0, SortBy: "created_at", SortDir: "desc",
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(out.Items) != 1 || out.Total != 1 || out.Items[0].AuthorLogin != "user123" {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestListListings_InvalidSortBy(t *testing.T) {
	repo := &fakeListRepo{}
	uc := usecase.NewListListings(repo)
	_, err := uc.Execute(context.Background(), usecase.ListListingsInput{
		SortBy: "unknown",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	var ae *apperror.AppError
	errors.As(err, &ae)
	if ae.Code != "invalid_sort_by" {
		t.Fatalf("wrong code: %s", ae.Code)
	}
}

func TestListListings_InvalidPriceRange(t *testing.T) {
	repo := &fakeListRepo{}
	uc := usecase.NewListListings(repo)
	minPrice := int64(100)
	maxPrice := int64(50)
	_, err := uc.Execute(context.Background(), usecase.ListListingsInput{
		PriceMin: &minPrice, PriceMax: &maxPrice,
	})
	if err == nil {
		t.Fatalf("expected range error")
	}
	var ae *apperror.AppError
	errors.As(err, &ae)
	if ae.Code != "invalid_price_range" {
		t.Fatalf("wrong code: %s", ae.Code)
	}
}
