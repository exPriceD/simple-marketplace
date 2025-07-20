package middleware

import (
	"net/http"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/platform/metrics"
)

func Metrics(m metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &statusWriter{ResponseWriter: w, status: 200}
			start := time.Now()
			next.ServeHTTP(sw, r)
			d := time.Since(start).Seconds()
			path := r.URL.Path
			m.IncHTTPRequests(r.Method, path, statusCodeLabel(sw.status))
			m.ObserveHTTPDuration(r.Method, path, d)
		})
	}
}

func statusCodeLabel(code int) string {
	return http.StatusText(code)
}
