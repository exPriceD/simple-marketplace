package user_test

import (
	userdomain "github.com/exPriceD/simple-marketplace/internal/domain/user"
	"github.com/exPriceD/simple-marketplace/internal/infrastructure/repository/user"
	"testing"
	"time"
)

func TestUserMapper_RoundTrip(t *testing.T) {
	loginVO, err := userdomain.NewLogin("Alice_01")
	if err != nil {
		t.Fatalf("login vo err: %v", err)
	}
	passHashVO, err := userdomain.NewPasswordHash("argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	if err != nil {
		t.Fatalf("pass hash vo err: %v", err)
	}
	created := time.Unix(1730000000, 0).UTC()

	orig := userdomain.RehydrateUser("abc123id", loginVO, passHashVO, created)

	row := userrepo.UserToRow(orig)
	if row.ID != "abc123id" || row.Login != "Alice_01" || row.PassHash != passHashVO.String() || !row.CreatedAt.Equal(created) {
		t.Fatalf("row mismatch %+v", row)
	}

	back, err := userrepo.RowToUser(row)
	if err != nil {
		t.Fatalf("RowToUser err: %v", err)
	}

	if back.ID() != orig.ID() ||
		back.Login().String() != orig.Login().String() ||
		back.PasswordHash() != orig.PasswordHash() ||
		!back.CreatedAt().Equal(orig.CreatedAt()) {
		t.Fatalf("round trip mismatch: got=%+v want=%+v", back, orig)
	}
}

func TestUserMapper_InvalidLogin(t *testing.T) {
	row := userrepo.UserRow{
		ID:        "id1",
		Login:     "bad login with space",
		PassHash:  "argon2id$v=19$t=1$m=65536$p=2$aaaaaaaaaaaaaaaaaaaaaa$bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		CreatedAt: time.Now().UTC(),
	}
	_, err := userrepo.RowToUser(row)
	if err == nil {
		t.Fatalf("expected error for invalid login")
	}
}

func TestUserMapper_InvalidPassHash(t *testing.T) {
	row := userrepo.UserRow{
		ID:        "id1",
		Login:     "Valid_123",
		PassHash:  "md5$insecure$hash",
		CreatedAt: time.Now().UTC(),
	}
	_, err := userrepo.RowToUser(row)
	if err == nil {
		t.Fatalf("expected error for invalid pass hash")
	}
}
