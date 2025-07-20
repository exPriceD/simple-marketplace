package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2idHasher реализует порт PasswordHasher.
type Argon2idHasher struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func NewArgon2idHasher(time, memoryKB uint32, threads uint8, keyLen uint32) *Argon2idHasher {
	return &Argon2idHasher{
		Time:    time,
		Memory:  memoryKB,
		Threads: threads,
		KeyLen:  keyLen,
	}
}

func (h *Argon2idHasher) Hash(plain string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt: %w", err)
	}
	key := argon2.IDKey([]byte(plain), salt, h.Time, h.Memory, h.Threads, h.KeyLen)
	return fmt.Sprintf("argon2id$v=19$t=%d$m=%d$p=%d$%s$%s",
		h.Time, h.Memory, h.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func (h *Argon2idHasher) Verify(hash, plain string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 7 {
		return false
	}
	if parts[0] != "argon2id" || parts[1] != "v=19" {
		return false
	}
	if !strings.HasPrefix(parts[2], "t=") ||
		!strings.HasPrefix(parts[3], "m=") ||
		!strings.HasPrefix(parts[4], "p=") {
		return false
	}
	tVal, err := strconv.ParseUint(strings.TrimPrefix(parts[2], "t="), 10, 32)
	if err != nil {
		return false
	}
	mVal, err := strconv.ParseUint(strings.TrimPrefix(parts[3], "m="), 10, 32)
	if err != nil {
		return false
	}
	pVal, err := strconv.ParseUint(strings.TrimPrefix(parts[4], "p="), 10, 8)
	if err != nil {
		return false
	}
	saltB64 := parts[5]
	keyB64 := parts[6]

	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false
	}
	keyStored, err := base64.RawStdEncoding.DecodeString(keyB64)
	if err != nil {
		return false
	}

	key := argon2.IDKey([]byte(plain), salt, uint32(tVal), uint32(mVal), uint8(pVal), uint32(len(keyStored)))
	return subtle.ConstantTimeCompare(keyStored, key) == 1
}
