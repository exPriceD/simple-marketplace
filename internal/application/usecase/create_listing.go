package usecase

import (
	"context"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
	"github.com/exPriceD/simple-marketplace/internal/application/dto/converters"
	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

type CreateListingInput struct {
	AuthorID    string
	Title       string
	Description string
	ImageURL    string
	Price       int64
}

type CreateListing struct {
	listings port.ListingRepository
	clock    port.Clock
}

func NewCreateListing(listings port.ListingRepository, clock port.Clock) *CreateListing {
	return &CreateListing{listings: listings, clock: clock}
}

func (uc *CreateListing) Execute(ctx context.Context, in CreateListingInput) (dto.ListingDTO, error) {
	op := "listing.create"
	if in.AuthorID == "" {
		return dto.ListingDTO{}, apperror.New(op, "unauthorized", apperror.KindAuth, nil)
	}
	titleVO, err := listing.NewTitle(in.Title)
	if err != nil {
		return dto.ListingDTO{}, apperror.New(op, "invalid_title", apperror.KindValidation, err)
	}
	descVO, err := listing.NewDescription(in.Description)
	if err != nil {
		return dto.ListingDTO{}, apperror.New(op, "invalid_description", apperror.KindValidation, err)
	}
	imgVO, err := listing.NewImageURL(in.ImageURL)
	if err != nil {
		return dto.ListingDTO{}, apperror.New(op, "invalid_image_url", apperror.KindValidation, err)
	}
	priceVO, err := listing.NewPrice(in.Price)
	if err != nil {
		return dto.ListingDTO{}, apperror.New(op, "invalid_price", apperror.KindValidation, err)
	}

	l := listing.NewListing(titleVO, descVO, imgVO, priceVO, in.AuthorID, uc.clock.Now())
	created, err := uc.listings.Create(ctx, l)
	if err != nil {
		return dto.ListingDTO{}, apperror.New(op, "persist_error", apperror.KindInternal, err)
	}
	return converters.ListingToDTO(created), nil
}
