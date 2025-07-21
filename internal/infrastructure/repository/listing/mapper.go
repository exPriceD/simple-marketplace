package listingrepo

import (
	"github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

// RowToListing преобразует ListingRow -> *listing.Listing.
func RowToListing(r ListingRow) (*listing.Listing, error) {
	titleVO, err := listing.NewTitle(r.Title)
	if err != nil {
		return nil, err
	}
	descVO, err := listing.NewDescription(r.Description)
	if err != nil {
		return nil, err
	}
	imgVO, err := listing.NewImageURL(r.ImageURL)
	if err != nil {
		return nil, err
	}
	priceVO, err := listing.NewPrice(r.Price)
	if err != nil {
		return nil, err
	}
	l := listing.RehydrateListing(r.ID, titleVO, descVO, imgVO, priceVO, r.AuthorID, r.AuthorLogin, r.CreatedAt)
	return l, nil
}

// ListingToRow преобразует доменную сущность -> ListingRow.
func ListingToRow(l *listing.Listing) ListingRow {
	return ListingRow{
		ID:          l.ID(),
		Title:       l.Title().String(),
		Description: l.Description(),
		ImageURL:    l.ImageURL(),
		Price:       l.Price(),
		AuthorID:    l.AuthorID(),
		AuthorLogin: l.AuthorLogin(),
		CreatedAt:   l.CreatedAt(),
	}
}
