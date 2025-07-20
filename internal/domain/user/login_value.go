package user

import (
	"regexp"

	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
)

var loginRegex = regexp.MustCompile(`^[A-Za-z0-9_]{3,32}$`)

type Login struct {
	value string
}

func NewLogin(raw string) (Login, error) {
	if !loginRegex.MatchString(raw) {
		return Login{}, shared.ErrInvalidLogin
	}
	return Login{value: raw}, nil
}

func (l Login) String() string { return l.value }
