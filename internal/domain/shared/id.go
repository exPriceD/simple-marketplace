package shared

import (
	"crypto/rand"
	"encoding/hex"
)

// IDGenerator генерирует 128-битный случайный ID в hex (32 символа).
type IDGenerator interface {
	NewID() string
}

type RandomIDGenerator struct{}

func (RandomIDGenerator) NewID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
