package listing_test

import (
	"strings"
	"testing"

	listingdomain "github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

func TestTitle_Success(t *testing.T) {
	_, err := listingdomain.NewTitle("Нормальный заголовок 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTitle_Fail(t *testing.T) {
	if _, err := listingdomain.NewTitle(""); err == nil {
		t.Fatalf("expected empty title error")
	}
	long := strings.Repeat("x", 121)
	if _, err := listingdomain.NewTitle(long); err == nil {
		t.Fatalf("expected too long error")
	}
}
