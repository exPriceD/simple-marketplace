package port

import (
	"context"

	"github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

type ListingRepository interface {
	Create(ctx context.Context, l *listing.Listing) (*listing.Listing, error)
	List(ctx context.Context, f ListingFilter) ([]*listing.Listing, error)
	Count(ctx context.Context, f ListingFilter) (int64, error)
}

type ListingFilter struct {
	Limit    int
	Offset   int
	SortBy   string
	SortDir  string
	PriceMin *int64
	PriceMax *int64
}
