package listing

import "github.com/exPriceD/simple-marketplace/internal/domain/shared"

type Title struct {
	value string
}

func NewTitle(raw string) (Title, error) {
	if len(raw) == 0 || len([]rune(raw)) > 120 {
		return Title{}, shared.ErrInvalidTitle
	}
	return Title{value: raw}, nil
}

func (t Title) String() string { return t.value }
