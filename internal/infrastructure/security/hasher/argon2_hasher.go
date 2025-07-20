package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// Argon2idHasher реализует порт PasswordHasher.
type Argon2idHasher struct {
	Time    uint32
	Memory  uint32 // в KB (обычно MB * 1024)
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
	var (
		t, m            uint32
		p               uint8
		saltB64, keyB64 string
	)
	_, err := fmt.Sscanf(hash, "argon2id$v=19$t=%d$m=%d$p=%d$%s$%s", &t, &m, &p, &saltB64, &keyB64)
	if err != nil {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return false
	}
	keyStored, err := base64.RawStdEncoding.DecodeString(keyB64)
	if err != nil {
		return false
	}
	key := argon2.IDKey([]byte(plain), salt, t, m, p, uint32(len(keyStored)))
	return subtle.ConstantTimeCompare(keyStored, key) == 1
}
