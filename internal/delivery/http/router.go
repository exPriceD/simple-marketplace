package httpdelivery

import (
	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"net/http"

	"github.com/exPriceD/simple-marketplace/internal/delivery/http/handler"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/middleware"
)

func NewRouter(authH *handler.AuthHandler, listingH *handler.ListingHandler, tp tokenParser) http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth
	mux.HandleFunc("/auth/register", method("POST", authH.Register))
	mux.HandleFunc("/auth/login", method("POST", authH.Login))

	// Listings
	mux.HandleFunc("/listings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listingH.List(w, r)
		case http.MethodPost:
			listingH.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	var h http.Handler = mux
	h = middleware.AuthContext(tp)(h)
	h = middleware.RequestID(h)
	h = middleware.Recovery(h)
	h = middleware.Logging(h)

	return h
}

type tokenParser interface {
	Parse(token string) (port.Claims, error)
}

// method — helper ограничения метода.
func method(method string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	}
}
