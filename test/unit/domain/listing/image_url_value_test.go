package listing_test

import (
	"testing"

	listingdomain "github.com/exPriceD/simple-marketplace/internal/domain/listing"
)

func TestImageURL_Success(t *testing.T) {
	_, err := listingdomain.NewImageURL("https://example.com/a.png")
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestImageURL_Fail(t *testing.T) {
	cases := []string{
		"", "ftp://x", "http://", "https://", "not-a-url",
	}
	for _, c := range cases {
		if _, err := listingdomain.NewImageURL(c); err == nil {
			t.Fatalf("expected error for %q", c)
		}
	}
}
