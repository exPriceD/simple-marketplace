package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
)

type ctxUserID struct{}
type ctxUserLogin struct{}

type tokenParser interface {
	Parse(token string) (port.Claims, error)
}

func AuthContext(tp tokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if h != "" {
				parts := strings.SplitN(h, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					if claims, err := tp.Parse(parts[1]); err == nil {
						ctx := context.WithValue(r.Context(), ctxUserID{}, claims.UserID)
						ctx = context.WithValue(ctx, ctxUserLogin{}, claims.Login)
						r = r.WithContext(ctx)
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxUserID{}).(string)
	return v, ok
}

func UserLogin(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxUserLogin{}).(string)
	return v, ok
}
