package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: 200}
		start := time.Now()
		next.ServeHTTP(sw, r)
		log.Printf("method=%s path=%s status=%d dur_ms=%d rid=%s ua=%q",
			r.Method, r.URL.String(), sw.status,
			time.Since(start).Milliseconds(), RequestIDValue(r.Context()),
			r.UserAgent(),
		)
	})
}
