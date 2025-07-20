package listing

import "github.com/exPriceD/simple-marketplace/internal/domain/shared"

type Price struct {
	value int64
}

func NewPrice(v int64) (Price, error) {
	if v < 0 || v > 1_000_000_000 {
		return Price{}, shared.ErrInvalidPrice
	}
	return Price{value: v}, nil
}

func (p Price) Int64() int64 { return p.value }
