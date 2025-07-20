package listing_test

import (
	"testing"

	listingdomain "github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

func TestPrice_Success(t *testing.T) {
	_, err := listingdomain.NewPrice(0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	_, err = listingdomain.NewPrice(1_000_000_000)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestPrice_Fail(t *testing.T) {
	if _, err := listingdomain.NewPrice(-1); err == nil {
		t.Fatalf("expected negative error")
	}
	if _, err := listingdomain.NewPrice(1_000_000_001); err == nil {
		t.Fatalf("expected upper bound error")
	}
}
