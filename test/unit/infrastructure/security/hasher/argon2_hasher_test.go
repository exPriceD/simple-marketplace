package hasher_test

import (
	"testing"

	"github.com/exPriceD/simple-marketplace/internal/infrastructure/security/hasher"
)

func TestArgon2Hasher_HashVerify(t *testing.T) {
	h := hasher.NewArgon2idHasher(1, 64*1024, 2, 32)
	hash, err := h.Hash("Password123!")
	if err != nil {
		t.Fatalf("hash err: %v", err)
	}
	if !h.Verify(hash, "Password123!") {
		t.Fatalf("should verify")
	}
	if h.Verify(hash, "Wrong") {
		t.Fatalf("should not verify wrong password")
	}
}

func TestArgon2Hasher_BadFormat(t *testing.T) {
	h := hasher.NewArgon2idHasher(1, 64*1024, 2, 32)
	if h.Verify("invalid-format", "x") {
		t.Fatalf("expected false for invalid format")
	}
}
