package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type ctxKeyRid struct{}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			var b [8]byte
			_, _ = rand.Read(b[:])
			rid = hex.EncodeToString(b[:])
		}
		ctx := context.WithValue(r.Context(), ctxKeyRid{}, rid)
		w.Header().Set("X-Request-ID", rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDValue(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRid{}).(string); ok {
		return v
	}
	return ""
}
