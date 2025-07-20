package listingrepo_test

import (
	"testing"
	"time"

	listingdomain "github.com/exPriceD/simple-marketplace/internal/domain/listing"
	listingrepo "github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/listing"
)

func TestListingMapper_RoundTrip(t *testing.T) {
	titleVO, err := listingdomain.NewTitle("Ps 5 Slim")
	if err != nil {
		t.Fatalf("title vo err: %v", err)
	}
	descVO, err := listingdomain.NewDescription("Новая, комплект полный, в подарок 1000 игр")
	if err != nil {
		t.Fatalf("desc vo err: %v", err)
	}
	imgVO, err := listingdomain.NewImageURL("https://example.com/i1.jpg")
	if err != nil {
		t.Fatalf("img vo err: %v", err)
	}
	priceVO, err := listingdomain.NewPrice(19990)
	if err != nil {
		t.Fatalf("price vo err: %v", err)
	}
	created := time.Unix(1730001000, 0).UTC()

	orig := listingdomain.RehydrateListing(42, titleVO, descVO, imgVO, priceVO, "user123", created)

	row := listingrepo.ListingToRow(orig)
	if row.ID != 42 || row.Title != titleVO.String() || row.Description != descVO.String() ||
		row.ImageURL != imgVO.String() || row.Price != priceVO.Int64() ||
		row.AuthorID != "user123" || !row.CreatedAt.Equal(created) {
		t.Fatalf("row mismatch: %+v", row)
	}

	back, err := listingrepo.RowToListing(row)
	if err != nil {
		t.Fatalf("RowToListing err: %v", err)
	}

	if back.ID() != orig.ID() ||
		back.Title().String() != orig.Title().String() ||
		back.Description() != orig.Description() ||
		back.ImageURL() != orig.ImageURL() ||
		back.Price() != orig.Price() ||
		back.AuthorID() != orig.AuthorID() ||
		!back.CreatedAt().Equal(orig.CreatedAt()) {
		t.Fatalf("round trip mismatch: got=%+v want=%+v", back, orig)
	}
}

func TestListingMapper_InvalidTitle(t *testing.T) {
	row := listingrepo.ListingRow{
		ID:          1,
		Title:       "",
		Description: "ok",
		ImageURL:    "https://example.com/img.jpg",
		Price:       100,
		AuthorID:    "u1",
		CreatedAt:   time.Now().UTC(),
	}
	_, err := listingrepo.RowToListing(row)
	if err == nil {
		t.Fatalf("expected error for invalid title")
	}
}

func TestListingMapper_InvalidPrice(t *testing.T) {
	row := listingrepo.ListingRow{
		ID:          1,
		Title:       "Valid",
		Description: "ok",
		ImageURL:    "https://example.com/img.jpg",
		Price:       -5,
		AuthorID:    "u1",
		CreatedAt:   time.Now().UTC(),
	}
	_, err := listingrepo.RowToListing(row)
	if err == nil {
		t.Fatalf("expected error for invalid price")
	}
}

func TestListingMapper_InvalidImageURL(t *testing.T) {
	row := listingrepo.ListingRow{
		ID:          1,
		Title:       "Valid",
		Description: "ok",
		ImageURL:    "ftp://bad",
		Price:       100,
		AuthorID:    "u1",
		CreatedAt:   time.Now().UTC(),
	}
	_, err := listingrepo.RowToListing(row)
	if err == nil {
		t.Fatalf("expected error for invalid image url")
	}
}
