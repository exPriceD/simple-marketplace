package converters

import (
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
	"github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

func ListingToDTO(l *listing.Listing) dto.ListingDTO {
	return dto.ListingDTO{
		ID:          l.ID(),
		Title:       l.Title().String(),
		Description: l.Description(),
		ImageURL:    l.ImageURL(),
		Price:       l.Price(),
		AuthorID:    l.AuthorID(),
		CreatedAt:   l.CreatedAt().UTC().Format(time.RFC3339),
	}
}
