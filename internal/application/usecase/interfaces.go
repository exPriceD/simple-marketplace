package usecase

import (
	"context"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
)

type RegisterUserUseCase interface {
	Execute(ctx context.Context, in RegisterUserInput) (dto.UserDTO, error)
}

type LoginUserUseCase interface {
	Execute(ctx context.Context, in LoginUserInput) (dto.AuthDTO, error)
}

type CreateListingUseCase interface {
	Execute(ctx context.Context, in CreateListingInput) (dto.ListingDTO, error)
}

type ListListingsUseCase interface {
	Execute(ctx context.Context, in ListListingsInput) (dto.ListingFeedDTO, error)
}
