package middleware

import (
	"net/http"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/platform/log"
)

func Logging(l log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w, status: 200}
			start := time.Now()
			next.ServeHTTP(sw, r)
			l.Info(r.Context(), "http_request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sw.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"rid", RequestIDValue(r.Context()),
			)
		})
	}
}
