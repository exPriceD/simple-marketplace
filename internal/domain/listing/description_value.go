package listing

import (
	"unicode/utf8"

	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
)

type Description struct {
	value string
}

func NewDescription(raw string) (Description, error) {
	if l := utf8.RuneCountInString(raw); l == 0 || l > 2000 {
		return Description{}, shared.ErrInvalidDescription
	}
	return Description{value: raw}, nil
}

func (d Description) String() string { return d.value }
