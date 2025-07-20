package system

import (
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
)

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }

var _ port.Clock = (*SystemClock)(nil)
