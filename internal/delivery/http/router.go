package httpdelivery

import (
	"net/http"
	"os"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/platform/log"
	"github.com/exPriceD/simple-marketplace/internal/platform/metrics"

	"github.com/exPriceD/simple-marketplace/internal/delivery/http/handler"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/middleware"
)

func NewRouter(
	authH *handler.AuthHandler,
	listingH *handler.ListingHandler,
	tp tokenParser,
	logger log.Logger,
	metr metrics.Metrics,
	metricsHandler http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Metrics
	mux.Handle("/metrics", metricsHandler)

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

	// Static
	fs := http.FileServer(http.Dir("frontend"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeFile(w, r, "frontend/index.html")
			return
		}

		if _, err := os.Stat("frontend" + r.URL.Path); os.IsNotExist(err) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			http.ServeFile(w, r, "frontend/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	var h http.Handler = mux
	h = middleware.AuthContext(tp)(h)
	h = middleware.Metrics(metr)(h)
	h = middleware.RequestID(h)
	h = middleware.Recovery(h)
	h = middleware.Logging(logger)(h)

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
