package user_test

import (
	"testing"

	userdomain "github.com/exPriceD/simple-marketplace/internal/domain/user"
)

func TestLogin_Success(t *testing.T) {
	valid := []string{"abc", "User_123", "A1_", "Z9Z9Z9Z9Z9Z9Z9Z9"}
	for _, v := range valid {
		if _, err := userdomain.NewLogin(v); err != nil {
			t.Fatalf("expected success for %q got %v", v, err)
		}
	}
}

func TestLogin_Fail(t *testing.T) {
	invalid := []string{"", "ab", "v   space", "русский", "toolooooooooooooooooooooooooooooooooooooooooooooooong"}
	for _, v := range invalid {
		if _, err := userdomain.NewLogin(v); err == nil {
			t.Fatalf("expected error for %q", v)
		}
	}
}
