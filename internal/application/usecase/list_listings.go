package usecase

import (
	"context"
	"strings"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
	"github.com/exPriceD/simple-marketplace/internal/application/dto/converters"
	apperror "github.com/exPriceD/simple-marketplace/internal/application/error"
	"github.com/exPriceD/simple-marketplace/internal/application/port"
)

type ListListingsInput struct {
	Limit         int
	Offset        int
	SortBy        string
	SortDir       string
	PriceMin      *int64
	PriceMax      *int64
	CurrentUserID string
}

type ListListings struct {
	listings port.ListingRepository
}

func NewListListings(listings port.ListingRepository) *ListListings {
	return &ListListings{listings: listings}
}

func (uc *ListListings) Execute(ctx context.Context, in ListListingsInput) (dto.ListingFeedDTO, error) {
	op := "listing.list"

	sortBy := in.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortBy != "created_at" && sortBy != "price" {
		return dto.ListingFeedDTO{}, apperror.New(op, "invalid_sort_by", apperror.KindValidation, nil)
	}
	sortDir := in.SortDir
	if sortDir == "" {
		sortDir = "desc"
	}
	sdLower := strings.ToLower(sortDir)
	if sdLower != "asc" && sdLower != "desc" {
		return dto.ListingFeedDTO{}, apperror.New(op, "invalid_sort_dir", apperror.KindValidation, nil)
	}
	if in.PriceMin != nil && in.PriceMax != nil && *in.PriceMin > *in.PriceMax {
		return dto.ListingFeedDTO{}, apperror.New(op, "invalid_price_range", apperror.KindValidation, nil)
	}

	filter := port.ListingFilter{
		Limit:    in.Limit,
		Offset:   in.Offset,
		SortBy:   sortBy,
		SortDir:  sdLower,
		PriceMin: in.PriceMin,
		PriceMax: in.PriceMax,
	}
	items, err := uc.listings.List(ctx, filter)
	if err != nil {
		return dto.ListingFeedDTO{}, apperror.New(op, "repo_error", apperror.KindInternal, err)
	}
	total, err := uc.listings.Count(ctx, filter)
	if err != nil {
		return dto.ListingFeedDTO{}, apperror.New(op, "count_error", apperror.KindInternal, err)
	}

	outItems := make([]dto.ListingDTO, 0, len(items))
	for _, l := range items {
		outItems = append(outItems, converters.ListingToDTO(l))
	}

	nextOffset := in.Offset + len(outItems)
	if nextOffset >= int(total) {
		nextOffset = -1
	}

	return dto.ListingFeedDTO{
		Items:      outItems,
		Total:      total,
		NextOffset: nextOffset,
	}, nil
}
