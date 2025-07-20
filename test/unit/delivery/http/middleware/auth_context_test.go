package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/exPriceD/simple-marketplace/internal/delivery/http/middleware"
)

type fakeTokenParser struct {
	out port.Claims
	err error
}

func (f fakeTokenParser) Parse(token string) (port.Claims, error) {
	return f.out, f.err
}

func TestAuthContext_Valid(t *testing.T) {
	parser := fakeTokenParser{out: port.Claims{
		UserID: "uid1", Login: "user1",
		Exp: time.Now().Add(time.Hour), Issued: time.Now(),
	}}
	h := middleware.AuthContext(parser)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := middleware.UserID(r.Context())
		if !ok || id != "uid1" {
			t.Fatalf("user id not set")
		}
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok123")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
}

func TestAuthContext_Invalid(t *testing.T) {
	parser := fakeTokenParser{err: assertErr("bad")}
	h := middleware.AuthContext(parser)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := middleware.UserID(r.Context()); ok {
			t.Fatalf("expected no user id")
		}
	}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer badtoken")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
