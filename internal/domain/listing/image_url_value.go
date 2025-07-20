package listing

import (
	"net/url"

	"github.com/exPriceD/simple-marketplace/internal/domain/shared"
)

type ImageURL struct {
	value string
}

func NewImageURL(raw string) (ImageURL, error) {
	if len(raw) == 0 || len(raw) > 512 {
		return ImageURL{}, shared.ErrInvalidImageURL
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ImageURL{}, shared.ErrInvalidImageURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ImageURL{}, shared.ErrInvalidImageURL
	}
	return ImageURL{value: raw}, nil
}

func (i ImageURL) String() string { return i.value }
