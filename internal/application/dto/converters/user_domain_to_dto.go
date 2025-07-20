package converters

import (
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/dto"
	"github.com/exPriceD/simple-marketplace/internal/domain/user"
)

func UserToDTO(u *user.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        u.ID(),
		Login:     u.Login().String(),
		CreatedAt: u.CreatedAt().UTC().Format(time.RFC3339),
	}
}
