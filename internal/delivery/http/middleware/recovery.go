package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/exPriceD/simple-marketplace/internal/delivery/http/response"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v\n%s", rec, debug.Stack())
				response.WriteInternalError(w, RequestIDValue(r.Context()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
